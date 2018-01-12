# go-clock-skew
## Description
This project is collect tcp timestamp by passive method and prepare for estimate [clock-skew](https://github.com/xiaoxuanzi/clock-skew).
## Usage
<pre><code>
Usage of ./go-clock-skew:<br>
  -e string<br>
    	device name (default "eth0")<br>
  -f string<br>
    	storage file (default "storage.csv")<br>
  -filter string<br>
    	bpFilter (default "tcp")<br>
  -h	help<br>
</pre></code>
## Example
<pre><code>
./go-clock-skew -filter "src host 10.10.89.144" -f 144.csv
</pre></code>
