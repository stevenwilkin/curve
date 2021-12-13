# Curve

Bitcoin futures curve from [Deribit](https://www.deribit.com/) as a JSON
webservice

## Building

	go build .

## Running

	./curve

Expiration date and annualised yield of each contract are then available via
HTTP:

	curl -s 0.0.0.0:8080 | jq '.'
	[
	  {
		"expiration": 1639728000000,
		"yield": 0.02933664033001339
	  },
	  {
		"expiration": 1640332800000,
		"yield": 0.0420198261426767
	  },
	  {
		"expiration": 1640937600000,
		"yield": 0.0580143048396088
	  },
	  {
		"expiration": 1643356800000,
		"yield": 0.06005654465966281
	  },
	  {
		"expiration": 1648195200000,
		"yield": 0.0718001618475805
	  },
	  {
		"expiration": 1656057600000,
		"yield": 0.08161726218932476
	  },
	  {
		"expiration": 1664524800000,
		"yield": 0.08491305394297999
	  }
	]

## Systemd service

Copy the service unit file to the configuration directory:

	cp curve.service /etc/systemd/system

Enable and start the service:

	systemctl enable curve
	systemctl start curve
