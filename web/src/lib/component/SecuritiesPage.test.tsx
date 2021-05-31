import React from 'react';
import {shallow} from 'enzyme';
import {RootStore} from '../store/RootStore';
import {newSecurityModel} from 'test/securityFactory';
import SecuritiesPage from './SecuritiesPage';
import Table from './table/Table';
import {IColumn} from './table/Column';
import TopAppBar from './TopAppBar';
import {SecurityModel} from '../model/SecurityModel';
import {ObservableMap} from 'mobx';
import {HideZero, Shares} from '../formats';

describe('SecuritiesPage', () => {
    const rootStore = new RootStore();
    const {securityStore} = rootStore;
    const security = newSecurityModel({shares: 96, firstAcquired: '1999-01-28', costBasis: 678.90, dividends: 345.67});

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue(rootStore);
    });
    it('displays app bar and table of securities', () => {
        const securities = [security];
        jest.spyOn(securityStore, 'securities', 'get').mockReturnValue(securities);

        const component = shallow(<SecuritiesPage />);

        expect(component.find(TopAppBar)).toHaveProp('title', 'Securities');
        expect(component.find(Table)).toHaveProp('data', securities);
    });
    describe('securities table', () => {
        beforeAll(() => {
            securityStore['securitiesById'] = new ObservableMap([
                [security.id, security],
            ]);
        });

        const columnTests: {input: SecurityModel, key: string, value: React.ReactNode, name?: string}[] = [
            {input: security, key: 'name', value: security.name},
            {input: security, key: 'type', value: security.type},
            {input: security, key: 'symbol', value: security.symbol},
            {input: security, key: 'shares', value: <Shares>{security.shares}</Shares>},
            {input: security, key: 'firstAcquired', value: security.firstAcquired},
            {input: security, key: 'costBasis', value: <HideZero>{security.costBasis}</HideZero>},
            {input: security, key: 'dividends', value: <HideZero>{security.dividends}</HideZero>},
            {input: security, key: 'transactionCount', value: security.transactionCount},
        ];

        columnTests.forEach(({name, input, key, value}) => {
            it(name ?? `displays ${key}`, () => {
                const component = shallow(<SecuritiesPage />);
                const columns = component.find(Table).prop<IColumn<SecurityModel>[]>('columns');
                const column = columns.find((column) => column.key === `security.${key}`);

                expect(column?.render(input)).toEqual(value);
            });
        });
        it('does not display zero shares', () => {
            const component = shallow(<SecuritiesPage />);
            const columns = component.find(Table).prop<IColumn<SecurityModel>[]>('columns');
            const column = columns.find((column) => column.key === 'security.shares');

            expect(column?.render(newSecurityModel())).toBeNull();
        });
        it('does not display cost basis for zero shares', () => {
            const component = shallow(<SecuritiesPage />);
            const columns = component.find(Table).prop<IColumn<SecurityModel>[]>('columns');
            const column = columns.find((column) => column.key === 'security.costBasis');

            expect(column?.render(newSecurityModel({costBasis: 0.000001}))).toBeNull();
        });
    });
});
