# Apollo

Apollo is a handy console-based way of organizing all of your media titles and keeping track of 
what you are currently watching, reading, and playing. It started out as a project for my own usage,
but I have since decided to put it out there on GitHub for all to enjoy.

Screenshot
----------

![ScreenShot](https://user-images.githubusercontent.com/3271352/38456471-61bbee10-3a5b-11e8-8e1c-330a0dafffd5.png)

Features
--------

- Create lists for movies, TV series, video games, and books.
- Sort the lists by year released, status, or ranking.
- Rank each entry on a scale of 1 to 6.
- Create collections of entries.
- Set each entry as being active, passive, or inactive.
- Print any list to a plain text file.

Commands
--------

- `/help` displays the help file.
- `/quit` quits the application.
- `/open <tab>` opens a specific tab.
- `/close` closes the current tab.
- `/set <config> <value>` sets an option to a specified value. 
- `/config` displays the current configuration.

Key-bindings
------------

- <kbd>Ctrl</kbd> <kbd>c</kbd> quits the application.
- <kbd>Alt</kbd> <kbd>[num]</kbd> switches to the num-th tab.
- <kbd>Enter</kbd> sends the current command, or toggles the input.
- <kbd>1</kbd> switches to the 'passive' view.
- <kbd>2</kbd> switches to the 'active' view.
- <kbd>3</kbd> switches to the 'inactive' view.
- <kbd>4</kbd> switches to the 'all' view.
- <kbd>s</kbd> sorts the entries.
- <kbd>D</kbd> deletes the current entry.
- <kbd>e</kbd> edits the current entry.
- <kbd>r</kbd> toggles ranking.
- <kbd>a</kbd> toggles the current entry's state.
- <kbd>[/]</kbd> changes the rating of the current entry.
- <kbd>Left/Right</kbd> changes the episodes of the current entry.
- <kbd>p</kbd> prints the current view to a file.

Installation
------------

This whole project was done using [Golang](https://golang.org/doc/install).

Once Go is installed properly, fetch this repository.

    go get github.com/finiks/apollo

Next, move to the repository source of the project and compile the application.

    cd $GOPATH/src/github/finiks/apollo
    go install

Lastly, you can run Apollo through its binary file.

    cd $GOPATH/bin
    ./apollo

At the first launch, `configuration.json` and `database.json` will be created and stored
in `~/.config/apollo/`.

Development State
-----------------

Apollo is now completely fulfilling my needs. That being said, if someone were to want to 
contribute to this project, please feel free to either message me for more information or 
to directly send in a pull request.
