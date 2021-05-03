import React from 'react';
import {shallow} from 'enzyme';
import TxDetail from './TxDetail';
import Category from './Category';
import {Currency, Shares} from 'src/lib/formats';
import Group from './Group';
import {newDetail} from 'src/test/detailFactory';
import {TextField} from '@material-ui/core';
import CategoryInput from './CategoryInput';
import GroupInput from './GroupInput';
import IconInput from '../IconInput';

const detail = newDetail({transactionGroupId: 456});

describe('TxDetail', () => {
    it('shows amount, category and group', () => {
        const component = shallow(<TxDetail detail={detail} editField={false} />);

        expect(component.find(Currency)).toHaveProp('children', detail.amount);
        expect(component.find(Category)).toHaveProp('detail', detail);
        expect(component.find(Group)).toHaveProp('id', detail.transactionGroupId);
        expect(component.find('.shares')).not.toExist();
        expect(component.find('.memo')).not.toExist();
    });
    it('shows asset quantity', () => {
        const assetQuantity = 100.789;

        const component = shallow(<TxDetail detail={{...detail, assetQuantity}} editField={false} />);

        expect(component.find('.shares').find(Shares)).toHaveProp('children', assetQuantity);
    });
    it('shows memo', () => {
        const memo = 'note to self';

        const component = shallow(<TxDetail detail={{...detail, memo}} editField={false} />);

        expect(component.find('.memo')).toHaveText(memo);
    });
    it('edits amount', () => {
        const component = shallow(<TxDetail detail={detail} editField='amount' />);

        const input = component.find(TextField);

        expect(input).toHaveProp({type: 'number', value: detail.amount, required: true});
    });
    it('edits category', () => {
        const component = shallow(<TxDetail detail={detail} editField='category' />);

        const input = component.find(CategoryInput);

        expect(input).toHaveProp('detail', detail);
    });
    it('edits group', () => {
        const component = shallow(<TxDetail detail={detail} editField='group' />);

        const input = component.find(GroupInput);

        expect(input).toHaveProp('detail', detail);
    });
    it('edits shares', () => {
        const detail = newDetail({assetQuantity: 123});
        const component = shallow(<TxDetail detail={detail} editField='shares' />);

        const input = component.find(IconInput);

        expect(input).toHaveProp({type: 'number', value: detail.assetQuantity, icon: 'request_page'});
    });
    it('edits memo', () => {
        const detail = newDetail({memo: 'some notes'});
        const component = shallow(<TxDetail detail={detail} editField='memo' />);

        const input = component.find(IconInput);

        expect(input).toHaveProp({value: detail.memo, icon: 'notes'});
    });
});
