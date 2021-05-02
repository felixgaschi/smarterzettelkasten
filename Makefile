.PHONY: 

install: 
	go build -o ~/smarterzettelkasten/szk
	echo "" >> ~/.bash_profile
	echo "# The next line updates PATH for Smarterzettelkasten" >> ~/.bash_profile
	printf "%s\n" $$'PATH = $$PATH:~/smarterzettelkasten/' >> ~/.bash_profile
	[ -f ~/.zshrc ] && echo "" >> ~/.zshrc
	[ -f ~/.zshrc ] && echo "# The next line updates PATH for Smarterzettelkasten" >> ~/.zshrc
	[ -f ~/.zshrc ] && printf "%s\n" $$'export PATH=$$PATH:~/smarterzettelkasten/' >> ~/.zshrc
