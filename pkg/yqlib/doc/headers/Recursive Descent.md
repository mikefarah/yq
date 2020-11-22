This operator recursively matches all children nodes given of a particular element, including that node itself. This is most often used to apply a filter recursively against all matches, for instance to set the `style` of all nodes in a yaml doc:

```bash
yq eval '.. style= "flow"' file.yaml
```