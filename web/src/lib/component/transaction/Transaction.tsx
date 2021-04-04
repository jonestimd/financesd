import {Typography, Checkbox, TextField, TextFieldProps} from '@material-ui/core';
import classNames from 'classnames';
import {observer} from 'mobx-react-lite';
import React from 'react';
import TransactionModel from '../../model/TransactionModel';
import Memo from './Memo';
import Payee from './Payee';
import Security from './Security';
import TxDetail from './TxDetail';
import * as formats from '../../formats';
import {translate} from '../../i18n/localize';
import PayeeInput from './PayeeInput';
import SecurityInput from './PayeeInput';
import DateField from '../DateField';

interface IProps {
    tx: TransactionModel;
    selected: boolean;
    showSecurity?: boolean;
}

const transactionFields = ['date', 'ref', 'payee', 'security', 'memo'] as const;
export type TransactionField = typeof transactionFields[number];

const detailFields = ['amount', 'category', 'group', 'shares', 'memo'] as const;
export type DetailField = typeof detailFields[number];

const securityField = 3;

const Transaction: React.FC<IProps> = observer(({tx, selected, showSecurity}) => {
    const [itemIndex, setItemIndex] = React.useState(-1);
    const [field, setField] = React.useState(0);
    const inputRef = React.useRef<HTMLDivElement | null>(null);
    React.useEffect(() => {
        inputRef.current?.querySelector('input')?.focus();
    }, [inputRef, itemIndex, field]);
    const transactionField = selected && itemIndex < 0 && transactionFields[field];
    const detailField = selected && detailFields[field];
    const onKeyDown = (event: React.KeyboardEvent) => {
        const {key, shiftKey} = event;
        if (key === 'Tab') {
            event.preventDefault();
            event.stopPropagation();
            if (shiftKey) {
                const f = field - 1;
                if (f === securityField && !showSecurity) setField(securityField-1);
                else if (f < 0) {
                    setField(transactionFields.length-1);
                    if (itemIndex === -1) setItemIndex(tx.details.length);
                    else setItemIndex(itemIndex-1);
                }
                else setField(f);
            }
            else {
                const f = field + 1;
                if (f === securityField && !showSecurity) setField(securityField+1);
                else if (f === transactionFields.length) {
                    setField(0);
                    if (itemIndex > tx.details.length) setItemIndex(-1);
                    else setItemIndex(itemIndex + 1);
                }
                else setField(f);
            }
        }
    };
    const inputProps: Partial<TextFieldProps> = {
        variant: 'filled',
        color: 'primary',
        size: 'small',
        ref: inputRef,
        onKeyDown,
    };
    return (
        <Typography className={classNames('transaction', {selected})}>
            <div className='leading'>
                {transactionField === 'date'
                    ? <DateField label={translate('Date')} value={tx.date ?? ''} {...inputProps} />
                    : <span className='date'>{tx.date}</span>}
                {transactionField === 'ref'
                    ? <TextField label={translate('Ref #')} value={tx.referenceNumber ?? ''} {...inputProps} />
                    : tx.referenceNumber ? <span className='ref-number'>{tx.referenceNumber}</span> : null}
            </div>
            <div className='details'>
                {transactionField === 'payee' ? <PayeeInput transaction={tx} {...inputProps} /> : <Payee transaction={tx} />}
                {transactionField === 'security' ? <SecurityInput transaction={tx} {...inputProps} /> : <Security transaction={tx} />}
                {transactionField === 'memo'
                    ? <TextField label={translate('Description')} value={tx.memo ?? ''} {...inputProps} />
                    : <Memo text={tx.memo} />}
                {tx.details.map((detail, index) =>
                    <TxDetail key={detail.id} detail={detail} editField={index === itemIndex && detailField} {...inputProps} />)}
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
