const ownerAvatars = [
	'/owner-avatar-1.svg',
	'/owner-avatar-2.svg',
	'/owner-avatar-3.svg',
	'/owner-avatar-4.svg',
	'/owner-avatar-5.svg'
];

export function getOwnerAvatarSrc(ownerId: number | null | undefined): string {
	if (!ownerId || ownerId < 1) {
		return ownerAvatars[0];
	}
	if (ownerId === 11) {
		return '/clown.png';
	}
	return ownerAvatars[(ownerId - 1) % ownerAvatars.length];
}
