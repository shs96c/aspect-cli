to run the demo from head just:
$ ./start.sh 
from this repo and then it will create a video of everything that is in the demo.sh file




scripts for demos: asciinema rec -c "./demo-template.sh" test


asciinema rec -c "./demo.sh" test

asciinema play test


requirements you need to install:

asciinema : https://asciinema.org/docs/installation
pv : http://www.ivarch.com/programs/pv.shtml

maybe add this to the ci for releases so videos are automatically created for every release