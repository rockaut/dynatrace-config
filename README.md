# dynatrace-config
Currently it's just an idea and a place to chat it out.

## Idea

Have a small program which utilizes yaml and/or json to talk to the dynatrace api to configure the settings (environment, configuration and maybe even cluster) from a repository through CI/CD. Why would just the developers have fun and not us infra-guys too?

### Questions
Golang or Python

### What there have to be

#### Ansible Module

Possible with go and python. No shortcomings in any way.

#### k8s pod / controller

- Pod could run constantly but also short-lived
- it would fetch multiple/single file(s) from defined sources (git, https) and apply the configs to target
- go would fit better as theres a smaller image

...
