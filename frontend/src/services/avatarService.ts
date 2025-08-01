const COLORS = [
  '#ff6b6b', '#f06595', '#cc5de8', '#845ef7', '#5c7cfa',
  '#339af0', '#22b8cf', '#20c997', '#51cf66', '#fcc419',
  '#ff922b', '#ff6b6b'
];

/**
 * Generates a consistent background color for a user's avatar based on their username.
 * @param username The username to generate a color for.
 * @returns A hex color code.
 */
export const generateAvatarColor = (username: string): string => {
  if (!username) {
    return '#ccc'; // Default color
  }
  let hash = 0;
  for (let i = 0; i < username.length; i++) {
    hash = username.charCodeAt(i) + ((hash << 5) - hash);
  }
  const index = Math.abs(hash % COLORS.length);
  return COLORS[index];
};
