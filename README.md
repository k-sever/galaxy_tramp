# Galaxy Tramp

In this game you are acting as a tramp in a galaxy full of black holes. 

Your goal is to find all the cells that are *not* black holes.

If you avoid a black hole, the number tells you how many of the 8 surrounding cells are a black hole. If it's blank, none of the surrounding tiles is a black hole.
Good luck!

## Build 
```shell
docker build . --target build
```
## Test
```shell
docker build . --target test --progress plain
```
## Run
```shell
docker build . --target game -t galaxy_tramp:latest
docker run -it galaxy_tramp:latest
```
You can specify one of 3 modes - easy(default), medium or hard. I.e.:
```shell
docker run -it galaxy_tramp:latest medium
```