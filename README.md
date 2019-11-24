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
- it would fetch multiple/single file(s) from defined sources (git, https) and apply the configs to target (the gitops way)
- go would fit better as theres a smaller image

...

## Experiments

### Simple go program which reads configuration json and calls the dynatrace api

#### Usage:
1. create a dynatrace trial account if you don't have one
2. create a dynatrace API-Token with Write Configuration in the dynatrace web console "Settings -> Integration -> Dynatrace API"
3. use the experiments/config/anomaly-detection-service.json or create a new file according to the https://www.dynatrace.com/support/help/extend-dynatrace/dynatrace-api/configuration-api/ . the configuration file must contain the API uri and the http method like in https://github.com/rockaut/dynatrace-config/blob/f26bf58b2d61a00b08f24a29ceb255d012954575/experiments/config/anomaly-detection-service.json#L2-L3 . The json body for the dynatrace API is in the payload key, like in https://github.com/rockaut/dynatrace-config/blob/f26bf58b2d61a00b08f24a29ceb255d012954575/experiments/config/anomaly-detection-service.json#L4
4. then execute (adapt to your use case)
```
export DYNATRACE_API_URL="https://yourenvironment.live.dynatrace.com/api/config"
export DYNATRACE_API_TOKEN="xxxxxx"
```

there are different command parameters available:

get current state of the configuration object:
```
go run experiments/dynatrace-config.go get experiments/config/anomaly-detection-service.json
```

show diff between current state and desired state of the configuration object:
```
go run experiments/dynatrace-config.go diff experiments/config/anomaly-detection-service.json
```

validate desired state:
```
go run experiments/dynatrace-config.go validate experiments/config/anomaly-detection-service.json
```

apply desired state:
```
go run experiments/dynatrace-config.go apply experiments/config/anomaly-detection-service.json
```

#### Findings:
- from my opinion it doesn't make sense to build a new API or configuration layer. In the configuration json files there needs to be the same payload as the dynatrace API server expects. Otherwise we always need to adapt our tool when the dynatrace API changes or will be expanded
- setting the dynatrace configuration is not always a PUT, sometimes also a POST. so we probably need to specify the http method in the configuration json file like in https://github.com/rockaut/dynatrace-config/blob/f26bf58b2d61a00b08f24a29ceb255d012954575/experiments/config/anomaly-detection-service.json#L3
- http returncode is depending on the api, so we also need to specify the expected success return code in the configuration json file, like in https://github.com/rockaut/dynatrace-config/blob/f26bf58b2d61a00b08f24a29ceb255d012954575/experiments/config/anomaly-detection-service.json#L36
- every (?) API has a 'validator' API which can validate the desired configuration. Not sure if that is really true for every API.
