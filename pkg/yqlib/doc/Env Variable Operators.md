
## Read boolean environment variable as a string
Running
```bash
myenv="true" yq eval --null-input 'strenv(myenv)'
```
will output
```yaml
12
```

## Read numeric environment variable as a string
Running
```bash
myenv="12" yq eval --null-input 'strenv(myenv)'
```
will output
```yaml
12
```

