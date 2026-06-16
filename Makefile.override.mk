# Local override of the generated `update-chart` target.
#
# The reusable sync-from-upstream workflow hard-codes `make update-chart`, and
# Makefile.gen.app.mk is generated (DO NOT EDIT). This file is loaded after it
# (include Makefile.*.mk sorts alphabetically: gen < override), so this recipe
# wins. Make prints an "overriding commands" warning — that is expected.
#
# Only this repo ships this file, so the shared workflow stays generic
# everywhere else; the override is the gate.

update-chart: check-env ## Sync chart with upstream, then run local sync/sync.sh.
	@echo "====> $@"
	./sync/sync.sh
	$(MAKE) update-deps
