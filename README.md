# pomodoro

Simple command-line productivity tool that implements the [Pomodoro technique](https://en.wikipedia.org/wiki/Pomodoro_Technique).

## Usage

```
pomodoro
> #########################
Pomodoro! 3 cycles today, rest for 5 minutes
```

```
pomodoro
> #########################
Pomodoro! 4 cycles today, take a longer break
```

It also pushes desktop notifications when the cycle times out. Check out the [requirements](https://github.com/0xAX/notificator) for your environment.

## Display your productivity on a terminal

```
pomodoro -graph
2017-10-31: ## 2
2017-11-01: ####### 7
2017-12-02: ##### 5
2017-12-03: ######## 8
```

Data is saved in JSON format in `pomodoro.dat` file in the same directory where the command runs.
