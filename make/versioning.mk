# MIT License
# Copyright (c) 2025 Toni Liesche
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.

tag-git-%: set-version-%
	git tag $(run.build.version) -m "Release $(run.build.version)"
	git push origin $(run.build.version)

set-version-release: set-additional-versions
	$(eval run.build.version := ${build.version.major}.${build.version.minor}.${build.version.bugfix})

set-version-rc: set-version-release
	$(eval run.build.version.rc := ${run.build.version}-rc)
	$(eval run.build.version.rc.minor := ${run.build.version.minor}-rc)
	$(eval run.build.version.rc.major := ${run.build.version.major}-rc)
	$(eval run.build.version := ${run.build.version}-rc${build.version.candidate})

set-version-patch: set-additional-versions
	$(eval run.build.version := ${build.version.major}.${build.version.minor}.${build.version.bugfix}.${build.version.patch})

set-additional-versions:
	$(eval run.build.version.major := ${build.version.major})
	$(eval run.build.version.minor := ${build.version.major}.${build.version.minor})

set-commit:
	$(eval run.commit := $(shell git rev-parse --short HEAD))

increase-%: update-% write-properties
	@echo updated build.properties file

update-major:
	@$(eval build.version.major := $(shell echo $$(($(build.version.major) + 1))))
	@$(eval build.version.minor := 0)
	@$(eval build.version.bugfix := 0)
	@$(eval build.version.candidate := 1)
	@$(eval build.version.patch := 1)
	@echo new major version: ${build.version.major}

update-minor:
	@$(eval build.version.minor := $(shell echo $$(($(build.version.minor) + 1))))
	@$(eval build.version.bugfix := 0)
	@$(eval build.version.candidate := 1)
	@$(eval build.version.patch := 1)
	@echo new minor version: ${build.version.minor}

update-bugfix:
	@$(eval build.version.bugfix := $(shell echo $$(($(build.version.bugfix) + 1))))
	@$(eval build.version.candidate := 1)
	@$(eval build.version.patch := 1)
	@echo new bugfix version: ${build.version.bugfix}

update-rc:
	@$(eval build.version.candidate := $(shell echo $$(($(build.version.candidate) + 1))))
	@echo new rc version: ${build.version.candidate}

update-patch:
	@$(eval build.version.patch := $(shell echo $$(($(build.version.patch) + 1))))
	@echo new patch version: ${build.version.patch}

write-properties:
	@echo "build.version.major=${build.version.major}" > build.properties.tmp
	@echo "build.version.minor=${build.version.minor}" >> build.properties.tmp
	@echo "build.version.bugfix=${build.version.bugfix}" >> build.properties.tmp
	@echo "build.version.candidate=${build.version.candidate}" >> build.properties.tmp
	@echo "build.version.patch=${build.version.patch}" >> build.properties.tmp
	@rm build.properties
	@mv build.properties.tmp build.properties
