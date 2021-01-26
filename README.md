# Pyrolytic graphite temperature experiment

## Rational
It has been found in an
[experiment](https://www.nextbigfuture.com/2020/10/tiny-energy-harvested-from-brownian-motion-could-replace-low-power-batteries.html)
that graphene vibrates thus producing usable energy.
Pyrolytic graphite is composed of many layers of graphene, so it should be warmer than the
surrounding environment.

## Experiment Setup
The experiment consists of one [k type thermocouple](https://en.wikipedia.org/wiki/Thermocouple#Type_K) attached to
[pyrolytic graphite](https://en.wikipedia.org/wiki/Pyrolytic_carbon) in a thermally
insulating container. Another thermocouple is outside of the container. Both thermocouples are
attached to a thermometer.

![setup](setup.png?raw=true)

## Experimental Results
The experiment measured the temperature difference between the two thermocouples averaged over a period
of about 80 minutes. On average the difference was found to be ~0.8F. The error in the measurement was
found to be less than 0.1F. The control experiment found the difference to be ~0.5F. This indicates that
the pyrolytic graphite is producing heat. The experiment was conducted at a room temperature of ~70F.

## Materials
* [pyrolytic graphite](https://unitednuclear.com/index.php?main_page=product_info&cPath=16_17_69&products_id=527)
* [meter](https://www.fluke.com/en-us/product/temperature-measurement/ir-thermometers/fluke-54-ii)
* [thermocouple](https://www.fluke.com/en-us/product/accessories/probes/fluke-80pk-1)

## Potential error
Pyrolytic graphite is conductive, so this could impact the thermocouple, but this would probably be noticable
in the meter readings.

## Data

### Pyrolytic graphite experiment - log1.csv
* average=0.764198
* corr=0.334671

![log1.csv](log1.png?raw=true)

### Pyrolytic graphite experiment with heat shrink tubing - log4.csv
* average=0.575000
* corr=-0.064983

![log4.csv](log4.png?raw=true)

### Pyrolytic graphite experiment with heat shrink tubing - 8 hours - night - log5.csv
* average=0.423647
* corr=0.825082

![log5.csv](log5.png?raw=true)

### Pyrolytic graphite experiment with heat shrink tubing - 8 hours - day - log6.csv
* average=0.437475
* corr=0.885628

![log6.csv](log6.png?raw=true)

### Calibration - log2.csv
* average=0.007059
* corr=0.993329

![log2.csv](log2.png?raw=true)

### Control - log3.csv
* average=0.460714
* corr=0.863758

![log3.csv](log3.png?raw=true)

### Control - 8 hours - night - log7.csv
* average=0.438677
* corr=0.574128

![log7.csv](log7.png?raw=true)

### Thermos Control - 8 hours - day  - log8.csv
* average=0.047695
* corr=0.971811

![log8.csv](log8.png?raw=true)

