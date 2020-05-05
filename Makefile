SHELL := bash
loc=$(echo ./checkpoints/$(date '+%d-%m-%Y')/)

checkpoint:
	mkdir $(value loc) && cp *json $(value loc)