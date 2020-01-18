# webhookcicd
 ### CICD automation tool to test build and push the containers to AWS ECR registry.  


 ```
$> cicdserver set tracking http://github.com/grapetechadmin/something
    OK
 
$> cicd --branch build --ver dev
    OK
    
$> cicd show
   version : ver
   current branch : master
   current build no : dev-x
   tracking : http://github.com/grapetechadmin/something
```
