import { redirect } from "next/navigation";

// / has no content of its own: send users to their todo list
// (/todos sends them to /login if they are not logged in)
export default function Home() {
  redirect("/todos");
}
