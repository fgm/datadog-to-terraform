# Jastify: JSON to Terraform 

A converter for Datadog dashboards and monitors.

## Installation (optional)

```
go install github.com/fgm/jastify@latest
```

If you do not wish to install the command you can also just run it like:

```
go run github.com/fgm/jastify@latest
```

## Convert DataDog JSON to Terraform in 3 steps

1. Create your monitor or dashboard in the Datadog UI and copy the JSON.
2. Two possibilities:
    - Using a file:
        - Save JSON to a file, say `dashboard.json`
        - Run `jastify dashboard.json`
        - Your resource will be named from the input file without the `.json` extension.
   - Using standard input
     - Run `jastify` in the terminal
     - Paste the JSON, terminating input with CTRL-D (CTRL-Z on Windows)
     - Since this does not allow you to provide a name, your resource
       will be named `dashboard_1` or `monitor_1`.
3. (optional) Jastify will output raw generated terraform code. 
Passing that through `terraform fmt` is recommended for legibility, e.g:
```
jastify dashboard.json | terraform fmt > dashboard.tf
```

## Contributing

We'd love for folks to contribute! Feel free to add your own ideas or take a look at the issues for inspiration.
The [Contributing Guide](CONTRIBUTING.md) explains development setup and the release process.

## Credits

This CLI converter is new code written based on the JavaScript
[datadog-to-terraform](https://github.com/laurmurclar/datadog-to-terraform)
browser extension, originally by Laura Murphy-Clarkin,
and its 2022 extensions in [PR#62](https://github.com/laurmurclar/datadog-to-terraform/pull/62)
published under the MIT License.
