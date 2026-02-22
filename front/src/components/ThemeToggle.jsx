import { useTheme } from '../ThemeContext'

/**
 * ThemeToggle - Minimal Neo-Operational theme switcher
 * Square button with sharp edges, no circular background
 */
export default function ThemeToggle() {
  const { theme, toggleTheme } = useTheme()

  return (
    <button
      onClick={toggleTheme}
      className="theme-toggle"
      aria-label="Toggle theme"
    >
      {theme === 'light' ? '◐' : '◑'}
    </button>
  )
}
