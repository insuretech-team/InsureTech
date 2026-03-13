/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				border: 'hsl(var(--color-border))',
				input: 'hsl(var(--color-input))',
				ring: 'hsl(var(--color-ring))',
				background: 'hsl(var(--color-background))',
				foreground: 'hsl(var(--color-foreground))',
				primary: {
					DEFAULT: 'hsl(var(--color-primary))',
					foreground: 'hsl(var(--color-primary-foreground))'
				},
				secondary: {
					DEFAULT: 'hsl(var(--color-secondary))',
					foreground: 'hsl(var(--color-secondary-foreground))'
				},
				destructive: {
					DEFAULT: 'hsl(var(--color-destructive))',
					foreground: 'hsl(var(--color-destructive-foreground))'
				},
				muted: {
					DEFAULT: 'hsl(var(--color-muted))',
					foreground: 'hsl(var(--color-muted-foreground))'
				},
				accent: {
					DEFAULT: 'hsl(var(--color-accent))',
					foreground: 'hsl(var(--color-accent-foreground))'
				},
				popover: {
					DEFAULT: 'hsl(var(--color-popover))',
					foreground: 'hsl(var(--color-popover-foreground))'
				},
				card: {
					DEFAULT: 'hsl(var(--color-card))',
					foreground: 'hsl(var(--color-card-foreground))',
					blue: 'hsl(var(--color-card-blue))',
					'blue-foreground': 'hsl(var(--color-card-blue-foreground))',
					purple: 'hsl(var(--color-card-purple))',
					'purple-foreground': 'hsl(var(--color-card-purple-foreground))',
					green: 'hsl(var(--color-card-green))',
					'green-foreground': 'hsl(var(--color-card-green-foreground))',
					orange: 'hsl(var(--color-card-orange))',
					'orange-foreground': 'hsl(var(--color-card-orange-foreground))'
				}
			},
			borderRadius: {
				lg: 'var(--radius)',
				md: 'calc(var(--radius) - 2px)',
				sm: 'calc(var(--radius) - 4px)'
			}
		}
	}
};
