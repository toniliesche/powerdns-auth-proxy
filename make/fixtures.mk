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

AUTH_PROXY_CMD=docker compose -f docker/docker-compose.yml -p powerdns-auth-proxy

setup-user: add-user add-role add-domain add-domain-role

add-user:
	@$(AUTH_PROXY_CMD) exec api powerdns-auth-proxy cli user add pdns-test-user password

list-users:
	@$(AUTH_PROXY_CMD) exec api powerdns-auth-proxy cli user list

delete-user:
	@$(AUTH_PROXY_CMD) exec api powerdns-auth-proxy cli user delete pdns-test-user

list-roles:
	@$(AUTH_PROXY_CMD) exec api powerdns-auth-proxy cli user-role list pdns-test-user

add-role:
	@$(AUTH_PROXY_CMD) exec api powerdns-auth-proxy cli user-role add pdns-test-user superadmin

remove-role:
	@$(AUTH_PROXY_CMD) exec api powerdns-auth-proxy cli user-role remove pdns-test-user superadmin

list-domains:
	@$(AUTH_PROXY_CMD) exec api powerdns-auth-proxy cli domain list

add-domain:
	@$(AUTH_PROXY_CMD) exec api powerdns-auth-proxy cli domain add toniliesche.de

delete-domain:
	@$(AUTH_PROXY_CMD) exec api powerdns-auth-proxy cli domain delete toniliesche.de

list-domain-roles:
	@$(AUTH_PROXY_CMD) exec api powerdns-auth-proxy cli domain-role list pdns-test-user toniliesche.de

add-domain-role:
	@$(AUTH_PROXY_CMD) exec api powerdns-auth-proxy cli domain-role add pdns-test-user toniliesche.de admin

remove-domain-role:
	@$(AUTH_PROXY_CMD) exec api powerdns-auth-proxy cli domain-role remove pdns-test-user toniliesche.de admin

list-api-keys:
	@$(AUTH_PROXY_CMD) exec api powerdns-auth-proxy cli api-key list pdns-test-user

add-api-key:
	@$(AUTH_PROXY_CMD) exec api powerdns-auth-proxy cli api-key add pdns-test-user

delete-api-key:
	@$(AUTH_PROXY_CMD) exec api powerdns-auth-proxy cli api-key delete pdns-test-user
