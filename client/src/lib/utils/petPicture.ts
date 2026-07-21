const PET_PICTURE_STORAGE_KEY = 'pet-picture-map-v1';
const PET_PICTURE_MAX_BYTES = 5 * 1024 * 1024;
const ALLOWED_PET_PICTURE_TYPES = new Set(['image/jpeg', 'image/png', 'image/webp']);

const defaultTypePictures: Record<string, string> = {
	camel: '/pet-camel.svg',
	snake: '/x-mark-red.png'
};

type PetPictureMap = Record<string, string>;

function loadPictureMap(): PetPictureMap {
	if (typeof localStorage === 'undefined') {
		return {};
	}
	try {
		const raw = localStorage.getItem(PET_PICTURE_STORAGE_KEY);
		if (!raw) return {};
		const parsed = JSON.parse(raw) as unknown;
		if (!parsed || typeof parsed !== 'object') return {};
		return parsed as PetPictureMap;
	} catch {
		return {};
	}
}

function persistPictureMap(map: PetPictureMap): void {
	if (typeof localStorage === 'undefined') {
		return;
	}
	localStorage.setItem(PET_PICTURE_STORAGE_KEY, JSON.stringify(map));
}

export function getPetPictureSrc(petId: number | null | undefined, typeName: string | null | undefined): string | null {
	if (petId && petId > 0) {
		const map = loadPictureMap();
		const customPicture = map[String(petId)];
		if (customPicture) {
			return customPicture;
		}
	}

	const typeKey = typeName?.trim().toLowerCase();
	if (!typeKey) return null;
	return defaultTypePictures[typeKey] ?? null;
}

export function hasCustomPetPicture(petId: number | null | undefined): boolean {
	if (!petId || petId < 1) return false;
	const map = loadPictureMap();
	return Boolean(map[String(petId)]);
}

export function removePetPicture(petId: number): void {
	const map = loadPictureMap();
	delete map[String(petId)];
	persistPictureMap(map);
}

export async function setPetPictureFromFile(petId: number, file: File): Promise<void> {
	if (!ALLOWED_PET_PICTURE_TYPES.has(file.type)) {
		throw new Error('Picture must be JPG, PNG, or WebP');
	}
	if (file.size <= 0) {
		throw new Error('Picture file is empty');
	}
	if (file.size > PET_PICTURE_MAX_BYTES) {
		throw new Error('Picture must be at most 5 MB');
	}

	const dataUrl = await readFileAsDataUrl(file);
	const map = loadPictureMap();
	map[String(petId)] = dataUrl;
	persistPictureMap(map);
}

function readFileAsDataUrl(file: File): Promise<string> {
	return new Promise((resolve, reject) => {
		const reader = new FileReader();
		reader.onload = () => resolve(String(reader.result ?? ''));
		reader.onerror = () => reject(new Error('Failed to read picture file'));
		reader.readAsDataURL(file);
	});
}
