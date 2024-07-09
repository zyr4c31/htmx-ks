air:
	templ generate
	tailwindcss -i ./src/input.css -o ./assets/css/styles.css
	air

templ:
	./scripts/templ.sh

wind:
	./scripts/tailwindcss.sh
