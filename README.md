# D7024E Kademlia 
![Go](https://github.com/viktorfrom/d7024e-kademlia/workflows/Go/badge.svg?branch=master)
[![GitHub license](https://img.shields.io/github/license/viktorfrom/d7024e-kademlia)](https://github.com/viktorfrom/d7024e-kademlia/blob/master/LICENSE)

Project designed and written in Go in conjunction with the D7024E Mobile and distributed computing systems course at Luleå University of Technology.

## Project description
The purpose of the project is to implement the Kademlia P2P distributed hash table network structure and simulate network communication between nodes.

## Requirements
* go 1.15+
* docker 19.03.12+

## Setup

### Golang 

### Linux
Below are the absolute minimum packages you will need for Linux. Names might vary depending on your distribution, you might need to install it manually if you can't find it using your distribution's package manager.
```
go 2:1.15-1
docker 1:19.03.12-2
```


## Build and Deploy
To build and deploy the Kademlia network run `scripts/deploy.sh`.


## While Running
List the different replica services
```
docker stack ps kadlab 
```

List the different running containers 
```
docker ps
```

Attach to a running container
```
docker attach "ContainerId"
```

## Test
To run the unit tests run `scripts/testcoverage.sh`.

## Authors
* Viktor From - vikfro-6@student.ltu.se - [viktorfrom](https://github.com/viktorfrom)
* Mark Hakansson - marhak-6@student.ltu.se - [markhakansson](https://github.com/markhakansson)
* Gustav Hansson - gushan-6@student.ltu.se - [97gushan](https://github.com/97gushan)

## License
Licensed under the MIT license. See [LICENSE](LICENSE) for details.
