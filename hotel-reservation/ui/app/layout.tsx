import Link from "next/link";
import { Providers } from "./providers";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <Link href="/">Home</Link>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
