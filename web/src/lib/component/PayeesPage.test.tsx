import React from 'react';
import {shallow} from 'enzyme';
import {RootStore} from '../store/RootStore';
import {newPayeeModel} from 'src/test/payeeFactory';
import PayeesPage from './PayeesPage';
import Table from './table/Table';
import {IColumn} from './table/Column';
import TopAppBar from './TopAppBar';
import {PayeeModel} from '../model/PayeeModel';
import {ObservableMap} from 'mobx';

describe('PayeesPage', () => {
    const rootStore = new RootStore();
    const {payeeStore} = rootStore;
    const payee = newPayeeModel();

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue(rootStore);
    });
    it('displays app bar and table of payees', () => {
        const payees = [payee];
        jest.spyOn(payeeStore, 'payees', 'get').mockReturnValue(payees);

        const component = shallow(<PayeesPage />);

        expect(component.find(TopAppBar)).toHaveProp('title', 'Payees');
        expect(component.find(Table)).toHaveProp('data', payees);
    });
    describe('payees table', () => {
        beforeAll(() => {
            payeeStore['payeesById'] = new ObservableMap([
                [payee.id, payee],
            ]);
        });

        const columnTests: {input: PayeeModel, key: string, value: React.ReactNode, name?: string}[] = [
            {input: payee, key: 'name', value: payee.name},
            {input: payee, key: 'transactionCount', value: payee.transactionCount},
        ];

        columnTests.forEach(({name, input, key, value}) => {
            it(name ?? `displays ${key}`, () => {
                const component = shallow(<PayeesPage />);
                const columns = component.find(Table).prop<IColumn<PayeeModel>[]>('columns');

                const column = columns.find((column) => column.key === `payee.${key}`);

                expect(column!.render(input)).toEqual(value);
            });
        });
    });
});
