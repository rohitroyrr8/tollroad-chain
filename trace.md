# Build steps

Built with Ignite v0.22.1.

* Run `docker build -f Dockerfile-exam . -t exam_i`.
* Run `docker run --rm -it -v $(pwd):/exam -w /exam exam_i ignite scaffold chain github.com/b9lab/toll-road`.
* Run `docker run --rm -it -v $(pwd):/exam -w /exam exam_i ignite scaffold single SystemInfo nextOperatorId:uint --no-message`
* Make it not nullable in genesis, fix compilation errors and tests.
* Run `docker run --rm -it -v $(pwd):/exam -w /exam exam_i ignite chain serve` and stop.
