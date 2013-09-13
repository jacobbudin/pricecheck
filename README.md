# pricecheck

Pricecheck is a command line tool to retrieve prices listed on the Web.

## Usage

Pricecheck requires two [YAML](http://www.yaml.org) configuration files: one for stores, and one for products. See the included YAML configuration files to understand their required format. The stores file requires an [XPath](https://developer.mozilla.org/en-US/docs/XPath) to locate the price in the DOM; these are similar to, but not compatible with, more familiar jQuery selectors.

(Note: Go is a compiled language. To run Pricecheck, you must first compile the application.)

From the command line, run like so:

	$ ./pricecheck -s stores.yaml -p products.yaml

## Quirks

- Cannot retrieve multiple prices from the same domain (e.g., display both the Nook and the hardcover price from Barnes & Noble)

## License

BSD 3-Clause (Revised) License
