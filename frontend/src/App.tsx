// App.tsx — the root component of your dashboard.
// Everything visible on the page starts here.

function App() {
  return (
    // The outermost div sets up a full-height flex column.
    // min-h-screen = "at least as tall as the viewport"
    // bg-background text-foreground = "use my dark theme colors"
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      {/* ─── TOP BAR ─────────────────────────────────────────── */}
      {/* border-b = bottom border. Uses --border token (subtle white line). */}
      <header className="border-b border-border">
        {/* max-w-7xl mx-auto = "center this, max 1280px wide" */}
        {/* px-6 py-4 = horizontal & vertical padding */}
        <div className="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
          <div className="flex items-center gap-2">
            {/* font-mono inherited from index.css (JetBrains Mono). */}
            <span className="text-lg font-semibold tracking-tight">
              Incident Platform
            </span>
          </div>
          {/* text-muted-foreground = dimmer text (the --muted-foreground token). */}
          <span className="text-sm text-muted-foreground">
            Monitor your services
          </span>
        </div>
      </header>

      {/* ─── MAIN CONTENT AREA ───────────────────────────────── */}
      {/* flex-1 = "grow to fill remaining vertical space" */}
      <main className="flex-1 max-w-7xl mx-auto w-full px-6 py-8">
        <h1 className="text-2xl font-semibold mb-6">Dashboard</h1>

        {/* Placeholder grid — we'll build real components next. */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="border border-border rounded-lg p-6">
            <p className="text-sm text-muted-foreground">Total Services</p>
            <p className="text-3xl font-semibold mt-2">—</p>
          </div>
          <div className="border border-border rounded-lg p-6">
            <p className="text-sm text-muted-foreground">Healthy</p>
            <p className="text-3xl font-semibold mt-2 text-emerald-500">—</p>
          </div>
          <div className="border border-border rounded-lg p-6">
            <p className="text-sm text-muted-foreground">Open Incidents</p>
            <p className="text-3xl font-semibold mt-2 text-red-500">—</p>
          </div>
        </div>
      </main>
    </div>
  )
}

export default App
