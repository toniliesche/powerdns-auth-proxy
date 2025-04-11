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

include make/docker_compose.mk
include make/docker_database.mk

create-network:
	@if ! docker network ls --format '{{.Name}}' | grep -q "^pdns-network$$"; then \
		docker network create pdns-network; \
		echo "Network pdns-network created."; \
	else \
		echo "Network pdns-network already exists."; \
	fi
