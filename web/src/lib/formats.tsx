import React from 'react';
import classNames from 'classnames';

export function numberClass(value: number, classes?: string) {
    return classNames(classes, 'number', {negative: value < 0});
}

export const shares = new Intl.NumberFormat(navigator.language, {style: 'decimal', minimumFractionDigits: 1, maximumFractionDigits: 6});
export const currency = new Intl.NumberFormat(navigator.language, {style: 'currency', currency: 'USD'});

export const Currency: React.FC<{children: number}> = ({children}) => {
    return <em className={numberClass(children)}>{currency.format(children)}</em>;
};

export const HideZero: React.FC<{children?: number}> = ({children}) => children ? <Currency>{children}</Currency> : null;

export const Shares: React.FC<{children: number}> = ({children}) => {
    return <em className={numberClass(children)}>{shares.format(children)}</em>;
};

const dateFormat = /^\d{4}-(0?[1-9]|1[0-2])-(0?[1-9]|[1-2][0-9]|3[01])$/;
const normalizeDate = (value: string) => {
    const parts = value.split('-');
    return `${parts[0].padStart(4, '0')}-${parts[1].padStart(2, '0')}-${parts[2].padStart(2, '0')}`;
};

export function isDate(text: string) {
    return dateFormat.test(text);
}

export function parseDate(text: string) {
    if (isDate(text)) return new Date(normalizeDate(text));
    return undefined;
}
