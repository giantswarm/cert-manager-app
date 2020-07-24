#!/bin/bash

DO_BACKUP=0
UPDATE_RESOURCES=0
UPDATE_SECRETS=0
UPDATE_INGRESSES=0
REMOVE_DEPRECATIONS=0

print_help() {
  echo -e "Usage:

  ./migrate-v090-to-v200.sh (--backup | --remove-deprecations | --update-ingresses --update-secrets [--update-cluster])

This script is designed to be used in three separate stages:

1. Back up existing resources.
2. Update existing resources with the new annotations and labels.
  1. This stage is explicitly a dry run.
3. Remove the deprecated annotations and labels.

Note: this script operates across **all namespaces**.

# Backing up resources

To back up all affected resources, pass the '--backup' flag. This will back up the following resources:

- Secret (of type 'kubernetes.io/tls', with deprecated labels/annotations)
- Ingress (where '.spec.tls' is set)
- Issuer
- ClusterIssuer
- Certificate
- CertificateRequest

Note: Files will be written to a subdirectory 'kubernetes-resources-backup 'alongside this script.
Note: '--update-cluster' and '--remove-deprecations' are mutually exclusive.

# Update resources

Resources which can be updated are of types Secret and Ingress.

To update affected resources, pass the flags '--update-ingresses' and '--update-secrets'. These can be provided
either together or independently. By default, this mode is a dry run and will not update the resources in your
cluster. To write the changes back, pass the '--update-cluster' flag.

Once complete, affected resources will have both old _and_ new annotations/labels - this allows you to ensure
that reconcilliation is functioning correctly whilst also providing a rollback option.

# Remove deprecations

To remove deprecated annotations/labels, pass the '--remove-deprecations' flag. This will remove deprecated
annotations and labels from resources of types Secret and Ingress.

Note: --update-cluster and --remove-deprecations are mutually exclusive."
}

do_backup() {
  BACKUP_DIR="kubernetes-resources-backup"
  declare RESOURCES=( 'ingress' 'secret' 'issuer' 'clusterissuer' 'certificate' 'certificaterequest' )

  # exit if the backup dir already exists
  if [[ -d ${BACKUP_DIR} ]]; then
    log_err "backup dir already exists at $(pwd)/${BACKUP_DIR}, aborting"
    return 1
  fi

  # make backup directory structure
  for resource in ${RESOURCES[@]}; do
    mkdir -p ${BACKUP_DIR}/${resource}
  done

  log_info "backup dir created at $(pwd)/${BACKUP_DIR}"
  log_info "backing up secrets to $(pwd)/${BACKUP_DIR}/secret/"

  # loop over namespaces
  for ns in $(kubectl get ns -o custom-columns=NAME:.metadata.name --no-headers=true); do
    resource_counter=0
    # loop over TLS secrets in namespace
    for secret in $(kubectl get secrets -n ${ns} -o custom-columns=NAME:.metadata.name --no-headers=true --field-selector type="kubernetes.io/tls"); do
      # only back up secrets with old API group set
      if [[ $(kubectl get secret ${secret} -n ${ns} -o json | grep "certmanager.k8s.io" | wc -l) -gt 0 ]]; then
        resource_counter=$((resource_counter+1))
        kubectl get secret ${secret} -n ${ns} -o json \
		| jq 'del(.metadata.resourceVersion,.metadata.uid,.status) | .metadata.creationTimestamp=null' \
		> ${BACKUP_DIR}/secret/${ns}_${secret}.json
      fi
    done
    log_info "${resource_counter} secret(s) backed up in namespace: ${ns}"
    resource_counter=0
  done

  log_info "backing up HTTPS ingresses to $(pwd)/${BACKUP_DIR}/ingress/"

  # loop over namespaces
  for ns in $(kubectl get ns -o custom-columns=NAME:.metadata.name --no-headers=true); do
    resource_counter=0
    # loop over ingresses in namespace
    for ingress in $(kubectl get ingress -n ${ns} -o custom-columns=NAME:.metadata.name --no-headers=true); do
      # only back up ingresses with .spec.tls set
      if [[ $(kubectl get ingress ${ingress} -n ${ns} -o json | jq '.spec | index("tls")') != "null" ]]; then
        resource_counter=$((resource_counter+1))
        kubectl get ingress ${ingress} -n ${ns} -o json \
		| jq 'del(.metadata.resourceVersion,.metadata.uid,.status) | .metadata.creationTimestamp=null' \
		> ${BACKUP_DIR}/ingress/${ns}_${ingress}.json
      fi
    done
    log_info "${resource_counter} ingress(es) backed up in namespace: ${ns}"
    ingress_counter=0
  done

  # back up clusterissuers (cluster-scoped)
  log_info "backing up clusterissuers to $(pwd)/${BACKUP_DIR}/clusterissuer/"

  resource_counter=0
  for clusterissuer in $(kubectl get clusterissuers -o custom-columns=NAME:.metadata.name --no-headers=true); do
    resource_counter=$((resource_counter+1))
    kubectl get clusterissuer ${clusterissuer} -o json \
	    | jq 'del(.metadata.resourceVersion,.metadata.uid,.status) | .metadata.creationTimestamp=null' \
            > ${BACKUP_DIR}/clusterissuer/${resource}.json
  done
  log_info "${resource_counter} clusterissuer(s) backed up"
  resource_counter=0

  # back up namespace-scoped cert-manager resources
  for resourcetype in 'issuer' 'certificate' 'certificaterequest'; do
    log_info "backing up ${resourcetype}s to $(pwd)/${BACKUP_DIR}/${resourcetype}/"

    # loop over namespaces
    for ns in $(kubectl get ns -o custom-columns=NAME:.metadata.name --no-headers=true); do
      resource_counter=0
      # loop over the resources of resourcetype in namespace
      for resource in $(kubectl get ${resourcetype} -n ${ns} -o custom-columns=NAME:.metadata.name --no-headers=true); do
        resource_counter=$((resource_counter+1))
        kubectl get ${resourcetype} ${resource} -n ${ns} -o json \
                | jq 'del(.metadata.resourceVersion,.metadata.uid,.status) | .metadata.creationTimestamp=null' \
                > ${BACKUP_DIR}/${resourcetype}/${ns}_${resource}.json
      done
      log_info "${resource_counter} ${resourcetype}(s) backed up in namespace: ${ns}"
      resource_counter=0
    done
  done
}

update_resources() {
  RESOURCE_TYPE=${1}

  if [[ ! ${UPDATE_RESOURCES} -eq 1 ]]; then
    log_info "DRY RUN ENABLED"
  fi

  log_info "updating resources of type: ${RESOURCE_TYPE}"

  # loop over namespaces
  for ns in $(kubectl get ns -o custom-columns=NAME:.metadata.name --no-headers=true); do
    # loop over resources in namespace
    if [[ ${RESOURCE_TYPE} == "ingress" ]]; then
      selector_command="kubectl get ingress -n ${ns} -o custom-columns=NAME:.metadata.name --no-headers=true"
    elif [[ ${RESOURCE_TYPE} == "secret" ]]; then
      selector_command="kubectl get secret -n ${ns} -o custom-columns=NAME:.metadata.name --no-headers=true --field-selector type='kubernetes.io/tls'"
    fi
    for resource in $(eval ${selector_command}); do
      # only operate on resources with the old API group set
      RESOURCE_JSON=$(kubectl get ${RESOURCE_TYPE} ${resource} -n ${ns} -o json)
      if grep -q "certmanager.k8s.io" <<< ${RESOURCE_JSON}; then
        if [[ ${UPDATE_RESOURCES} -eq 1 ]]; then
          log_info "updating ${RESOURCE_TYPE} ${resource} in namespace ${ns}"
          RESOURCE_JSON=$(sed 's/certmanager.k8s.io/cert-manager.io/g' <<< ${RESOURCE_JSON})
          kubectl apply -f - <<< ${RESOURCE_JSON}
        else
          log_info "DRY RUN: updating ${RESOURCE_TYPE} ${resource} in namespace ${ns}"
        fi
      fi
      unset -v RESOURCE_JSON
    done
  done
}

remove_deprecations() {

  log_info "removing deprecated API group from secrets"

  # loop over namespaces
  for ns in $(kubectl get ns -o custom-columns=NAME:.metadata.name --no-headers=true); do
    # loop over secrets in namespace
    for secret in $(kubectl get secret -n ${ns} -o custom-columns=NAME:.metadata.name --no-headers=true --field-selector type="kubernetes.io/tls"); do
      # only operate on secrets with old API group set
      if [[ $(kubectl get secret ${secret} -n ${ns} -o json | jq -r .metadata.annotations | grep "certmanager.k8s.io" | wc -l) -gt 0 ]]; then
        log_info "Operating on secret '${secret}' in namespace '${ns}'"
        # loop over deprecated annotations
        for annotation in $(kubectl get secret ${secret} -n ${ns} -o json \
                | jq -r '.metadata.annotations | with_entries(select(.key|match("certmanager.k8s.io")))' \
                | jq -r '. | keys[]') ; do
          log_info "removing annotation '${annotation}'"
          annotation=$(sed 's|/|~1|g' <<< ${annotation})
          kubectl_cmd="kubectl patch secret ${secret} -n ${ns} --type=json -p='[{"op": "remove", "path": "/metadata/annotations/${annotation}"}]'"
          eval "${kubectl_cmd}"
        done
        # loop over deprecated labels
        for label in $(kubectl get secret ${secret} -n ${ns} -o json \
		| jq -r '.metadata.labels | with_entries(select(.key|match("certmanager.k8s.io")))' \
		| jq -r '. | keys[]') ; do
	  log_info "removing label '${label}'"
	  label=$(sed 's|/|~1|g' <<< ${label})
	  kubectl_cmd="kubectl patch secret ${secret} -n ${ns} --type=json -p='[{"op": "remove", "path": "/metadata/labels/${label}"}]'"
          eval "${kubectl_cmd}"
        done
      fi
    done
  done
}

log () {
  level=$1
  shift 1
  date -u +"%Y-%m-%dT%H:%M:%SZ" | tr -d '\n'
  echo " [${level}] $@"
}

log_info () {
  log "INFO" "$@"
}

log_err () {
  log "ERR" "$@"
}

parse_args () {

  if [[ $# -eq 0 ]]; then
    print_help
    exit 0
  fi

  while [[ $# -gt 0 ]]; do
    case $1 in
      -h|--help)
        print_help
        exit 0
        ;;
      --backup)
        DO_BACKUP=1
        shift
        ;;
      --update-cluster)
        UPDATE_RESOURCES=1
        shift
        ;;
      --remove-deprecations)
        REMOVE_DEPRECATIONS=1
        shift
        ;;
      --update-ingresses)
        UPDATE_INGRESSES=1
        shift
        ;;
      --update-secrets)
        UPDATE_SECRETS=1
        shift
        ;;
   esac
  done

  if [[ ${DO_BACKUP} -eq 1 ]]; then
    do_backup
    exit 0
  fi

  if ! which jq 1>/dev/null 2>&1 ; then
    log_err "jq is not installed - it is required for this script"
    exit 1
  fi

  if [[ ${UPDATE_RESOURCES} -eq 1 ]] && [[ ${REMOVE_DEPRECATIONS} -eq 1 ]]; then
    log_err "'--update-resources' & '--remove-deprecations' cannot be used together"
    exit 1
  fi

  if [[ ${REMOVE_DEPRECATIONS} -eq 1 ]]; then
    remove_deprecations
  fi

  if [[ ${UPDATE_INGRESSES} -eq 1 ]]; then
    update_resources ingresses
  fi
  if [[ ${UPDATE_SECRETS} -eq 1 ]]; then
    update_resources secrets
  fi
}

parse_args $@
