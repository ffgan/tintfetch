# Security

This project is in development and has no stable supported release yet.
Security fixes target the current development branch. Do not assume historical
commits receive fixes.

Report vulnerabilities privately through the repository's GitHub Security
Advisory reporting interface if enabled. If it is unavailable, open an issue
asking the maintainer for a private reporting channel without including exploit
details or sensitive information. There is no guaranteed response SLA.

Native declarative configuration is parsed as data. `--shell-config` explicitly
executes trusted Bash code with the invoking user's permissions. External
engines retain their own execution and configuration behavior. Public-IP lookup
is optional and contacts the URL explicitly configured for that module.
