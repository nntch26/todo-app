import Link from "next/link";

import { Button } from "@/components/ui/Button";
import type { User } from "@/types/auth";

type NavbarProps = {
  user: User;
  onLogout: () => void;
};

// Navbar: top bar on pages that need login (app name, user email, log out)
export function Navbar({ user, onLogout }: NavbarProps) {
  return (
    <header className="border-b border-border bg-surface">
      <nav className="mx-auto flex max-w-3xl items-center justify-between gap-4 px-6 py-4">
        {/* App name */}
        <Link href="/todos" className="text-xl font-semibold">
          Todo<span className="text-accent">.</span>
        </Link>

        {/* User + log out */}
        <div className="flex items-center gap-4">
          {/* Hidden on small screens to save space */}
          <span className="hidden text-sm text-muted sm:inline">{user.email}</span>

          <Button variant="secondary" size="sm" onClick={onLogout}>
            Log out
          </Button>
        </div>
      </nav>
    </header>
  );
}
