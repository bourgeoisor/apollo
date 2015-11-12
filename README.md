# Apollo

Apollo is a handy console-based way of organizing all of your media titles and keeping track of 
what you are currently watching, reading, and playing. It started out as a project for my own usage,
but I have since decided to put it out there on Github for all to enjoy.

Screenshot
----------

![ScreenShot](https://cloud.githubusercontent.com/assets/3271352/11099708/4c142eb0-8883-11e5-923f-4e28d7b7b76e.png)

Features
--------

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
- <kbd>t</kbd> auto-tags the current entry.
- <kbd>r</kbd> toggles ranking.
- <kbd>a</kbd> toggles the current entry's state.
- <kbd>z/x</kbd> changes the rating of the current entry.
- <kbd>c/v</kbd> changes the episodes of the current entry.
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

All in all, Apollo as a CLI application is basically done, minus some little features that I might want
to implement later on. That being said, if someone were to want to contribute to this project, please 
feel free to either message me for more information or to directly send in a pull request. The next
step for Apollo will be to port the application as an Android and possibly iOS application.

APIs Used
---------

Various public APIs were used to fetch metadata from the media titles.

- [OMDb](http://omdbapi.com/)
- [Hummingbird](https://github.com/hummingbird-me/hummingbird/wiki/API-v1-Methods)
- [TheGamesDB](http://wiki.thegamesdb.net/index.php/API_Introduction)
- [Google Books](https://developers.google.com/books/docs/overview)
