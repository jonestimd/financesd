import React from 'react';
import {shallow} from 'enzyme';
import * as formats from './formats';
import {isDate, parseDate} from './formats';

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
    describe('HideZero', () => {
        it('does not display 0', () => {
            const result = shallow(<formats.HideZero>{0}</formats.HideZero>);

            expect(result).toBeEmptyRender();
        });
        it('displays non-zero value as currency', () => {
            const result = shallow(<formats.HideZero>{0.01}</formats.HideZero>);

            expect(result.find(formats.Currency)).toHaveProp('children', 0.01);
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
    describe('isDate', () => {
        it('returns true if date matches year-month-day', () => {
            expect(isDate('1234-1-1')).toBe(true);
            expect(isDate('2345-01-01')).toBe(true);
            expect(isDate('3456-12-31')).toBe(true);
        });
        it('returns false if year not 4 digits', () => {
            expect(isDate('345-11-02')).toBe(false);
            expect(isDate('34567-11-02')).toBe(false);
        });
        it('returns false if month > 12', () => {
            expect(isDate('3456-13-02')).toBe(false);
        });
        it('returns false if month is 0', () => {
            expect(isDate('3456-0-02')).toBe(false);
            expect(isDate('3456-00-02')).toBe(false);
        });
        it('returns false if day > 31', () => {
            expect(isDate('3456-12-32')).toBe(false);
        });
        it('returns false if day is 0', () => {
            expect(isDate('3456-12-0')).toBe(false);
            expect(isDate('3456-12-00')).toBe(false);
        });
    });
    describe('parseDate', () => {
        it('returns undefined for invalid date', () => {
            expect(parseDate('2021-06-0')).toBeUndefined();
            expect(parseDate('2021-0-01')).toBeUndefined();
            expect(parseDate('202-01-01')).toBeUndefined();
            expect(parseDate('20201-01-01')).toBeUndefined();
        });
        it('returns date', () => {
            expect(parseDate('2021-06-03')).toEqual(new Date('2021-06-03'));
            expect(parseDate('2021-6-03')).toEqual(new Date('2021-06-03'));
            expect(parseDate('2021-6-3')).toEqual(new Date('2021-06-03'));
        });
    });
});
