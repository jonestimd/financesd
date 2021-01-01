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

export const Shares: React.FC<{children: number}> = ({children}) => {
    return <em className={numberClass(children)}>{shares.format(children)}</em>;
};
