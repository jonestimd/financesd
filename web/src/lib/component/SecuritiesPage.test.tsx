import React from 'react';
import {shallow} from 'enzyme';
import {RootStore} from '../store/RootStore';
import {newSecurityModel} from 'src/test/securityFactory';
import SecuritiesPage from './SecuritiesPage';
import Table, {IColumn} from './table/Table';
import TopAppBar from './TopAppBar';
import {SecurityModel} from '../model/SecurityModel';
import {ObservableMap} from 'mobx';

describe('SecuritiesPage', () => {
    const rootStore = new RootStore();
    const {securityStore} = rootStore;
    const security = newSecurityModel();

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
            {input: security, key: 'transactionCount', value: security.transactionCount},
        ];

        columnTests.forEach(({name, input, key, value}) => {
            it(name ?? `displays ${key}`, () => {
                const component = shallow(<SecuritiesPage />);
                const columns = component.find(Table).prop<IColumn<SecurityModel>[]>('columns');

                const column = columns.find((column) => column.key === `security.${key}`);

                expect(column.render(input)).toEqual(value);
            });
        });
    });
});
