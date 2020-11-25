---
name: Feature request - V4
about: Suggest an idea for this project
title: ''
labels: [enhancement, v4]
assignees: ''

---

**Please describe your feature request.**
A clear and concise description of what the request is and what it would solve. 
Ex. I wish I could use yq to [...]

Please note that V3 will no longer have any enhancements.

**Describe the solution you'd like**
If we have data1.yml like:
(please keep to around 10 lines )

```yaml
country: Australia
```

And we run a command:

```bash
yq eval 'predictWeatherOf(.country)'
```

it could output

```yaml
temp: 32
```

**Describe alternatives you've considered**
A clear and concise description of any alternative solutions or features you've considered.

**Additional context**
Add any other context or screenshots about the feature request here.
