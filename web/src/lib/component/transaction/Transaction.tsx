import {Typography, Checkbox, TextFieldProps} from '@material-ui/core';
import classNames from 'classnames';
import {observer} from 'mobx-react-lite';
import React from 'react';
import TransactionModel from '../../model/TransactionModel';
import Memo from './Memo';
import Payee from './Payee';
import Security from './Security';
import TxDetail from './TxDetail';
import * as formats from '../../formats';
import PayeeInput from './PayeeInput';
import SecurityInput from './SecurityInput';
import DateInput from '../DateInput';
import IconInput from '../IconInput';

interface IProps {
    tx: TransactionModel;
    selected: boolean;
    showSecurity?: boolean;
    fieldIndex: number;
    setField: (field: number) => void;
}

const unselectedField = {transactionField: false as const, detailField: false as const, itemIndex: -1};

const Transaction: React.FC<IProps> = observer(({tx, selected, showSecurity, fieldIndex, setField}) => {
    const {transactionField, detailField, itemIndex} = selected ? tx.getField(fieldIndex, showSecurity) : unselectedField;
    const inputProps: Partial<TextFieldProps> = {
        variant: 'outlined',
        color: 'primary',
        size: 'small',
        onKeyDown: (event: React.KeyboardEvent<HTMLInputElement>) => {
            const {key, ctrlKey, shiftKey, altKey} = event;
            if (key === 'Tab' && !ctrlKey && !altKey) {
                event.stopPropagation();
                event.preventDefault();
                setField(tx.nextField(fieldIndex + (shiftKey ? -1 : 1), showSecurity));
            }
        },
    };
    return (
        <Typography className={classNames('transaction', {selected})}>
            <div className='leading'>
                {transactionField === 'date'
                    ? <DateInput value={tx.date ?? ''} {...inputProps} />
                    : <span className='date'>{tx.date}</span>}
                {transactionField === 'ref'
                    ? <IconInput value={tx.referenceNumber ?? ''} {...inputProps} icon='tag' />
                    : tx.referenceNumber ? <span data-type='ref'>{tx.referenceNumber}</span> : null}
            </div>
            <div className='details'>
                {transactionField === 'payee' ? <PayeeInput transaction={tx} {...inputProps} /> : <Payee transaction={tx} />}
                {transactionField === 'security' ? <SecurityInput transaction={tx} {...inputProps} /> : <Security transaction={tx} />}
                {transactionField === 'description'
                    ? <IconInput value={tx.memo ?? ''} {...inputProps} icon='notes' />
                    : <Memo text={tx.memo} />}
                {tx.details.map((detail, index) =>
                    <TxDetail key={detail.id} detail={detail} editField={index === itemIndex && detailField} showSecurity={showSecurity} {...inputProps} />)}
            </div>
            <div className='trailing'>
                <Checkbox disabled checked={tx.cleared} />
                <span className={formats.numberClass(tx.subtotal)}>{formats.currency.format(tx.subtotal)}</span>
                <span className={formats.numberClass(tx.balance)}>{formats.currency.format(tx.balance)}</span>
            </div>
        </Typography>
    );
});

export default Transaction;
