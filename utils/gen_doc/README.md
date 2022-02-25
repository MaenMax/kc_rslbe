gen_doc
===========
Generate RAML files out of the routes registered in cumulis_fe3/src/empowerthings.com/cumulis/cumulis_fe3/router.  

There will be as many RAML files generated as top nodes. For example, if we have the following top nodes: v3.0 and apps, then v3.0.raml and apps.raml files will be generated.  


### How to visualize the RAML files  

You will need an extra tool for that. They can be found on Internet. Just google.  

Personally, I'm using [raml-doc](https://github.com/nidi3/raml-doc) as follow:  

```bash

git clone https://github.com/nidi3/raml-doc

cd raml-doc

mvn package

# raml-doc doesn't seem to like ABSOLUTE path, so you will have to copy your RAML file
# in the local directory or somewhere local so that you can use RELATIVE path.
# I suppose I copy it in current directory (in raml-doc top directory).

java -jar raml-doc-standalone/target/raml-doc-standalone-0.8.3-SNAPSHOT.jar ./v3.0.raml

# This will generate a raml-doc directory with Javascript/CSS/HTML files and many others.
# Then you can vizualize by using Firefox or Chrome and open the index.html file from
# your computer.

```


