import React from 'react';
import {shallow} from 'enzyme';
import TransactionList from './TransactionList';
import {RootStore} from 'lib/store/RootStore';
import TransactionTableModel from 'lib/model/TransactionTableModel';
import ListViewport from '../scroll/ListViewport';
import {newTx, newTxModel} from 'test/transactionFactory';
import Memo from './Memo';
import {Checkbox} from '@material-ui/core';
import {newDetail} from 'test/detailFactory';
import TransactionModel from 'lib/model/TransactionModel';
import * as selectionHooks from '../scroll/listSelectionHooks';
import {mockListSelectionHook} from 'test/mockHooks';
import Transaction from './Transaction';

const Wrapper: React.FC<{children: React.ReactElement}> = ({children}) => <>{children}</>;

describe('TransactionList', () => {
    const {categoryStore, transactionStore, accountStore} = new RootStore();
    const detail = newDetail({amount: 123.45});
    const transaction = newTx({memo: 'transaction memo', details: [detail]});
    const transactionsModel = new TransactionTableModel([transaction], accountStore, categoryStore);

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue({transactionStore, accountStore});
        jest.spyOn(transactionStore, 'getTransactionsModel').mockReturnValue(transactionsModel);
    });
    it('displays list in ListViewPort', () => {
        const selection = mockListSelectionHook();

        const component = shallow(<TransactionList />);

        expect(component.find(ListViewport).props()).toEqual(expect.objectContaining({
            items: transactionsModel.transactions, selection,
        }));
    });
    it('displays prototype transaction', () => {
        const component = shallow(<TransactionList />);

        const prototype = component.find('TransactionPrototype').dive();
        expect(prototype.find('.leading .date')).toHaveText('0000-00-00');
        expect(prototype.find('.details').find(Memo)).toHaveProp('text', 'Prototype');
        expect(prototype.find('.trailing').find(Checkbox)).toExist();
        expect(prototype.find('.trailing .number').at(0)).toHaveText('0.00');
    });
    describe('selection', () => {
        it('resets field when row changes', () => {
            const useSelection = jest.spyOn(selectionHooks, 'useSelection');
            shallow(<TransactionList />);

            const getColumn = useSelection.mock.calls[0][0].getColumn!;

            expect(getColumn({row: Math.random(), column: Math.random()})).toEqual(0);
        });
    });
    describe('renderItem', () => {
        const renderTx = (tx: TransactionModel, selected: boolean) => {
            const renderItem = shallow(<TransactionList />).find(ListViewport).prop('renderItem');
            return shallow(<Wrapper>{renderItem(tx, -1, selected) as React.ReactElement}</Wrapper>);
        };

        const txModel = newTxModel();
        it('displays transaction', () => {
            const tx = renderTx(txModel, false);

            expect(tx.find(Transaction)).toHaveProp({tx: txModel, selected: false, fieldIndex: 0});
        });
        it('highlights selected transaction', () => {
            const tx = renderTx(txModel, true);

            expect(tx.find(Transaction)).toHaveProp('selected', true);
        });
        it('displays input for selected field', () => {
            mockListSelectionHook(0, 3);

            const tx = renderTx(txModel, true);

            expect(tx.find(Transaction)).toHaveProp('fieldIndex', 3);
        });
        it('updates column when field changes', () => {
            const selection = mockListSelectionHook(0, 3);
            const tx = renderTx(txModel, true);

            tx.find(Transaction).prop<(f: number) => void>('setField')(4);

            expect(selection.setCell).toBeCalledWith(0, 4);
        });
    });
});
