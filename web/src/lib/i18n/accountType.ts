import {Bundle, defaultKey, translate} from './localize';

const bundle: Bundle = {
    [defaultKey]: {
        BANK: "Bank",
        BROKERAGE: "Brokerage",
        CASH: "Cash",
        CREDIT: "Credit",
        LOAN: "Loan",
        _401K: "401(K)",
    }
};

export default function get(code: string) {
    return translate(code, bundle);
}