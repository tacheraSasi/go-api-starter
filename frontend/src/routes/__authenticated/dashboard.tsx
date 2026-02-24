import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/__authenticated/dashboard')({
  component: () => <div>Hello /authenticated/dashboard!</div>,
});
