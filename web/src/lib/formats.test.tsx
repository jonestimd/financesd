import React from 'react';
import {shallow} from 'enzyme';
import * as formats from './formats';

describe('formats', () => {
    describe('numberClass', () => {
        it('adds number to class names', () => {
            expect(formats.numberClass(0, 'other')).toEqual('other number');
        });
        it('adds negative to class names if value is < 0', () => {
            expect(formats.numberClass(-0.01, 'other')).toEqual('other number negative');
        });
    });
    describe('shares', () => {
        it('rounds to 6 digits', () => {
            expect(formats.shares.format(1.2345678)).toEqual('1.234568');
        });
    });
    describe('currency', () => {
        it('displays dollar amount', () => {
            expect(formats.currency.format(1.23)).toEqual('$1.23');
        });
    });
    describe('Currency', () => {
        it('displays children as dollar amount', () => {
            const result = shallow(<formats.Currency>{1.23}</formats.Currency>);

            expect(result.type()).toEqual('em');
            expect(result).toHaveClassName('number');
            expect(result).not.toHaveClassName('negative');
            expect(result).toHaveText('$1.23');
        });
        it('adds negative class if value < 0', () => {
            const result = shallow(<formats.Currency>{-1.23}</formats.Currency>);

            expect(result.type()).toEqual('em');
            expect(result).toHaveClassName('negative');
            expect(result).toHaveText('-$1.23');
        });
    });
    describe('Shares', () => {
        it('displays children as asset quantity', () => {
            const result = shallow(<formats.Shares>{1.2345678}</formats.Shares>);

            expect(result.type()).toEqual('em');
            expect(result).toHaveClassName('number');
            expect(result).not.toHaveClassName('negative');
            expect(result).toHaveText('1.234568');
        });
        it('adds negative class if value < 0', () => {
            const result = shallow(<formats.Shares>{-1.2345678}</formats.Shares>);

            expect(result.type()).toEqual('em');
            expect(result).toHaveClassName('negative');
            expect(result).toHaveText('-1.234568');
        });
    });
});
