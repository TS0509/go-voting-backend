// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// VotingCandidate is an auto generated low-level Go binding around an user-defined struct.
type VotingCandidate struct {
	Name      string
	VoteCount *big.Int
}

// VotingMetaData contains all meta data concerning the Voting contract.
var VotingMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"candidateIndex\",\"type\":\"uint256\"}],\"name\":\"Voted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"}],\"name\":\"addCandidate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"candidates\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"endVoting\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCandidates\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"internalType\":\"structVoting.Candidate[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"hasEnded\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"hasStarted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"hasVoted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"startVoting\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"candidateIndex\",\"type\":\"uint256\"}],\"name\":\"vote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// VotingABI is the input ABI used to generate the binding from.
// Deprecated: Use VotingMetaData.ABI instead.
var VotingABI = VotingMetaData.ABI

// Voting is an auto generated Go binding around an Ethereum contract.
type Voting struct {
	VotingCaller     // Read-only binding to the contract
	VotingTransactor // Write-only binding to the contract
	VotingFilterer   // Log filterer for contract events
}

// VotingCaller is an auto generated read-only Go binding around an Ethereum contract.
type VotingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VotingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VotingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VotingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VotingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VotingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VotingSession struct {
	Contract     *Voting           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VotingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VotingCallerSession struct {
	Contract *VotingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// VotingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VotingTransactorSession struct {
	Contract     *VotingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VotingRaw is an auto generated low-level Go binding around an Ethereum contract.
type VotingRaw struct {
	Contract *Voting // Generic contract binding to access the raw methods on
}

// VotingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VotingCallerRaw struct {
	Contract *VotingCaller // Generic read-only contract binding to access the raw methods on
}

// VotingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VotingTransactorRaw struct {
	Contract *VotingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVoting creates a new instance of Voting, bound to a specific deployed contract.
func NewVoting(address common.Address, backend bind.ContractBackend) (*Voting, error) {
	contract, err := bindVoting(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Voting{VotingCaller: VotingCaller{contract: contract}, VotingTransactor: VotingTransactor{contract: contract}, VotingFilterer: VotingFilterer{contract: contract}}, nil
}

// NewVotingCaller creates a new read-only instance of Voting, bound to a specific deployed contract.
func NewVotingCaller(address common.Address, caller bind.ContractCaller) (*VotingCaller, error) {
	contract, err := bindVoting(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VotingCaller{contract: contract}, nil
}

// NewVotingTransactor creates a new write-only instance of Voting, bound to a specific deployed contract.
func NewVotingTransactor(address common.Address, transactor bind.ContractTransactor) (*VotingTransactor, error) {
	contract, err := bindVoting(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VotingTransactor{contract: contract}, nil
}

// NewVotingFilterer creates a new log filterer instance of Voting, bound to a specific deployed contract.
func NewVotingFilterer(address common.Address, filterer bind.ContractFilterer) (*VotingFilterer, error) {
	contract, err := bindVoting(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VotingFilterer{contract: contract}, nil
}

// bindVoting binds a generic wrapper to an already deployed contract.
func bindVoting(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VotingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Voting *VotingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Voting.Contract.VotingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Voting *VotingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Voting.Contract.VotingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Voting *VotingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Voting.Contract.VotingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Voting *VotingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Voting.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Voting *VotingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Voting.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Voting *VotingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Voting.Contract.contract.Transact(opts, method, params...)
}

// Candidates is a free data retrieval call binding the contract method 0x3477ee2e.
//
// Solidity: function candidates(uint256 ) view returns(string name, uint256 voteCount)
func (_Voting *VotingCaller) Candidates(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Name      string
	VoteCount *big.Int
}, error) {
	var out []interface{}
	err := _Voting.contract.Call(opts, &out, "candidates", arg0)

	outstruct := new(struct {
		Name      string
		VoteCount *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Name = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.VoteCount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Candidates is a free data retrieval call binding the contract method 0x3477ee2e.
//
// Solidity: function candidates(uint256 ) view returns(string name, uint256 voteCount)
func (_Voting *VotingSession) Candidates(arg0 *big.Int) (struct {
	Name      string
	VoteCount *big.Int
}, error) {
	return _Voting.Contract.Candidates(&_Voting.CallOpts, arg0)
}

// Candidates is a free data retrieval call binding the contract method 0x3477ee2e.
//
// Solidity: function candidates(uint256 ) view returns(string name, uint256 voteCount)
func (_Voting *VotingCallerSession) Candidates(arg0 *big.Int) (struct {
	Name      string
	VoteCount *big.Int
}, error) {
	return _Voting.Contract.Candidates(&_Voting.CallOpts, arg0)
}

// GetCandidates is a free data retrieval call binding the contract method 0x06a49fce.
//
// Solidity: function getCandidates() view returns((string,uint256)[])
func (_Voting *VotingCaller) GetCandidates(opts *bind.CallOpts) ([]VotingCandidate, error) {
	var out []interface{}
	err := _Voting.contract.Call(opts, &out, "getCandidates")

	if err != nil {
		return *new([]VotingCandidate), err
	}

	out0 := *abi.ConvertType(out[0], new([]VotingCandidate)).(*[]VotingCandidate)

	return out0, err

}

// GetCandidates is a free data retrieval call binding the contract method 0x06a49fce.
//
// Solidity: function getCandidates() view returns((string,uint256)[])
func (_Voting *VotingSession) GetCandidates() ([]VotingCandidate, error) {
	return _Voting.Contract.GetCandidates(&_Voting.CallOpts)
}

// GetCandidates is a free data retrieval call binding the contract method 0x06a49fce.
//
// Solidity: function getCandidates() view returns((string,uint256)[])
func (_Voting *VotingCallerSession) GetCandidates() ([]VotingCandidate, error) {
	return _Voting.Contract.GetCandidates(&_Voting.CallOpts)
}

// HasEnded is a free data retrieval call binding the contract method 0xecb70fb7.
//
// Solidity: function hasEnded() view returns(bool)
func (_Voting *VotingCaller) HasEnded(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Voting.contract.Call(opts, &out, "hasEnded")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasEnded is a free data retrieval call binding the contract method 0xecb70fb7.
//
// Solidity: function hasEnded() view returns(bool)
func (_Voting *VotingSession) HasEnded() (bool, error) {
	return _Voting.Contract.HasEnded(&_Voting.CallOpts)
}

// HasEnded is a free data retrieval call binding the contract method 0xecb70fb7.
//
// Solidity: function hasEnded() view returns(bool)
func (_Voting *VotingCallerSession) HasEnded() (bool, error) {
	return _Voting.Contract.HasEnded(&_Voting.CallOpts)
}

// HasStarted is a free data retrieval call binding the contract method 0x44691f7e.
//
// Solidity: function hasStarted() view returns(bool)
func (_Voting *VotingCaller) HasStarted(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Voting.contract.Call(opts, &out, "hasStarted")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasStarted is a free data retrieval call binding the contract method 0x44691f7e.
//
// Solidity: function hasStarted() view returns(bool)
func (_Voting *VotingSession) HasStarted() (bool, error) {
	return _Voting.Contract.HasStarted(&_Voting.CallOpts)
}

// HasStarted is a free data retrieval call binding the contract method 0x44691f7e.
//
// Solidity: function hasStarted() view returns(bool)
func (_Voting *VotingCallerSession) HasStarted() (bool, error) {
	return _Voting.Contract.HasStarted(&_Voting.CallOpts)
}

// HasVoted is a free data retrieval call binding the contract method 0x09eef43e.
//
// Solidity: function hasVoted(address ) view returns(bool)
func (_Voting *VotingCaller) HasVoted(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Voting.contract.Call(opts, &out, "hasVoted", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasVoted is a free data retrieval call binding the contract method 0x09eef43e.
//
// Solidity: function hasVoted(address ) view returns(bool)
func (_Voting *VotingSession) HasVoted(arg0 common.Address) (bool, error) {
	return _Voting.Contract.HasVoted(&_Voting.CallOpts, arg0)
}

// HasVoted is a free data retrieval call binding the contract method 0x09eef43e.
//
// Solidity: function hasVoted(address ) view returns(bool)
func (_Voting *VotingCallerSession) HasVoted(arg0 common.Address) (bool, error) {
	return _Voting.Contract.HasVoted(&_Voting.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Voting *VotingCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Voting.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Voting *VotingSession) Owner() (common.Address, error) {
	return _Voting.Contract.Owner(&_Voting.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Voting *VotingCallerSession) Owner() (common.Address, error) {
	return _Voting.Contract.Owner(&_Voting.CallOpts)
}

// AddCandidate is a paid mutator transaction binding the contract method 0x462e91ec.
//
// Solidity: function addCandidate(string _name) returns()
func (_Voting *VotingTransactor) AddCandidate(opts *bind.TransactOpts, _name string) (*types.Transaction, error) {
	return _Voting.contract.Transact(opts, "addCandidate", _name)
}

// AddCandidate is a paid mutator transaction binding the contract method 0x462e91ec.
//
// Solidity: function addCandidate(string _name) returns()
func (_Voting *VotingSession) AddCandidate(_name string) (*types.Transaction, error) {
	return _Voting.Contract.AddCandidate(&_Voting.TransactOpts, _name)
}

// AddCandidate is a paid mutator transaction binding the contract method 0x462e91ec.
//
// Solidity: function addCandidate(string _name) returns()
func (_Voting *VotingTransactorSession) AddCandidate(_name string) (*types.Transaction, error) {
	return _Voting.Contract.AddCandidate(&_Voting.TransactOpts, _name)
}

// EndVoting is a paid mutator transaction binding the contract method 0xc3403ddf.
//
// Solidity: function endVoting() returns()
func (_Voting *VotingTransactor) EndVoting(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Voting.contract.Transact(opts, "endVoting")
}

// EndVoting is a paid mutator transaction binding the contract method 0xc3403ddf.
//
// Solidity: function endVoting() returns()
func (_Voting *VotingSession) EndVoting() (*types.Transaction, error) {
	return _Voting.Contract.EndVoting(&_Voting.TransactOpts)
}

// EndVoting is a paid mutator transaction binding the contract method 0xc3403ddf.
//
// Solidity: function endVoting() returns()
func (_Voting *VotingTransactorSession) EndVoting() (*types.Transaction, error) {
	return _Voting.Contract.EndVoting(&_Voting.TransactOpts)
}

// StartVoting is a paid mutator transaction binding the contract method 0x1ec6b60a.
//
// Solidity: function startVoting() returns()
func (_Voting *VotingTransactor) StartVoting(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Voting.contract.Transact(opts, "startVoting")
}

// StartVoting is a paid mutator transaction binding the contract method 0x1ec6b60a.
//
// Solidity: function startVoting() returns()
func (_Voting *VotingSession) StartVoting() (*types.Transaction, error) {
	return _Voting.Contract.StartVoting(&_Voting.TransactOpts)
}

// StartVoting is a paid mutator transaction binding the contract method 0x1ec6b60a.
//
// Solidity: function startVoting() returns()
func (_Voting *VotingTransactorSession) StartVoting() (*types.Transaction, error) {
	return _Voting.Contract.StartVoting(&_Voting.TransactOpts)
}

// Vote is a paid mutator transaction binding the contract method 0x0121b93f.
//
// Solidity: function vote(uint256 candidateIndex) returns()
func (_Voting *VotingTransactor) Vote(opts *bind.TransactOpts, candidateIndex *big.Int) (*types.Transaction, error) {
	return _Voting.contract.Transact(opts, "vote", candidateIndex)
}

// Vote is a paid mutator transaction binding the contract method 0x0121b93f.
//
// Solidity: function vote(uint256 candidateIndex) returns()
func (_Voting *VotingSession) Vote(candidateIndex *big.Int) (*types.Transaction, error) {
	return _Voting.Contract.Vote(&_Voting.TransactOpts, candidateIndex)
}

// Vote is a paid mutator transaction binding the contract method 0x0121b93f.
//
// Solidity: function vote(uint256 candidateIndex) returns()
func (_Voting *VotingTransactorSession) Vote(candidateIndex *big.Int) (*types.Transaction, error) {
	return _Voting.Contract.Vote(&_Voting.TransactOpts, candidateIndex)
}

// VotingVotedIterator is returned from FilterVoted and is used to iterate over the raw logs and unpacked data for Voted events raised by the Voting contract.
type VotingVotedIterator struct {
	Event *VotingVoted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VotingVotedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VotingVoted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VotingVoted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VotingVotedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VotingVotedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VotingVoted represents a Voted event raised by the Voting contract.
type VotingVoted struct {
	Voter          common.Address
	CandidateIndex *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterVoted is a free log retrieval operation binding the contract event 0x4d99b957a2bc29a30ebd96a7be8e68fe50a3c701db28a91436490b7d53870ca4.
//
// Solidity: event Voted(address indexed voter, uint256 indexed candidateIndex)
func (_Voting *VotingFilterer) FilterVoted(opts *bind.FilterOpts, voter []common.Address, candidateIndex []*big.Int) (*VotingVotedIterator, error) {

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}
	var candidateIndexRule []interface{}
	for _, candidateIndexItem := range candidateIndex {
		candidateIndexRule = append(candidateIndexRule, candidateIndexItem)
	}

	logs, sub, err := _Voting.contract.FilterLogs(opts, "Voted", voterRule, candidateIndexRule)
	if err != nil {
		return nil, err
	}
	return &VotingVotedIterator{contract: _Voting.contract, event: "Voted", logs: logs, sub: sub}, nil
}

// WatchVoted is a free log subscription operation binding the contract event 0x4d99b957a2bc29a30ebd96a7be8e68fe50a3c701db28a91436490b7d53870ca4.
//
// Solidity: event Voted(address indexed voter, uint256 indexed candidateIndex)
func (_Voting *VotingFilterer) WatchVoted(opts *bind.WatchOpts, sink chan<- *VotingVoted, voter []common.Address, candidateIndex []*big.Int) (event.Subscription, error) {

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}
	var candidateIndexRule []interface{}
	for _, candidateIndexItem := range candidateIndex {
		candidateIndexRule = append(candidateIndexRule, candidateIndexItem)
	}

	logs, sub, err := _Voting.contract.WatchLogs(opts, "Voted", voterRule, candidateIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VotingVoted)
				if err := _Voting.contract.UnpackLog(event, "Voted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVoted is a log parse operation binding the contract event 0x4d99b957a2bc29a30ebd96a7be8e68fe50a3c701db28a91436490b7d53870ca4.
//
// Solidity: event Voted(address indexed voter, uint256 indexed candidateIndex)
func (_Voting *VotingFilterer) ParseVoted(log types.Log) (*VotingVoted, error) {
	event := new(VotingVoted)
	if err := _Voting.contract.UnpackLog(event, "Voted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
