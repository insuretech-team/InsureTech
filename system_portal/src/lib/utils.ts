import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export type WithElementRef<Props, Ref = Element> = Props & {
	ref?: Ref | null;
};

export type WithoutChild<Props> = Props extends { child?: unknown }
	? Omit<Props, 'child'>
	: Props;

export type WithoutChildren<Props> = Props extends { children?: unknown }
	? Omit<Props, 'children'>
	: Props;

export type WithoutChildrenOrChild<Props> = WithoutChildren<WithoutChild<Props>>;
