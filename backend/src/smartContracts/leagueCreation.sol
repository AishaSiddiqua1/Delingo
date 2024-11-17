// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract LeagueCompetition {
    struct League {
        address creator;
        string name;
        uint256 startTime;
        uint256 endTime;
        bool isActive;
        uint256 entryFee;
        uint256 rewardPool; // Total reward pool in tokens
    }

    struct Participant {
        address user;
        uint256 points;
    }

    mapping(uint256 => League) public leagues;
    mapping(uint256 => Participant[]) public leagueParticipants;

    uint256 public leagueCount = 0;
    IERC20 public lingToken; // Reference to the LingToken contract
    address public treasury; // Address of the treasury holding LingTokens

    event LeagueCreated(uint256 leagueId, string name, address indexed creator);
    event JoinedLeague(uint256 leagueId, address indexed participant);
    event PointsUpdated(uint256 leagueId, address indexed participant, uint256 points);
    event RewardsDistributed(uint256 leagueId);

    constructor(address _lingToken, address _treasury) {
        lingToken = IERC20(_lingToken);
        treasury = _treasury;
    }

    // Create a new league
    function createLeague(
        string memory _name,
        uint256 _startTime,
        uint256 _endTime,
        uint256 _entryFee,
        uint256 _rewardPool
    ) public {
        require(_startTime < _endTime, "Start time must be before end time");
        require(
            lingToken.allowance(treasury, address(this)) >= _rewardPool,
            "Insufficient reward pool allowance"
        );

        leagueCount++;
        leagues[leagueCount] = League({
            creator: msg.sender,
            name: _name,
            startTime: _startTime,
            endTime: _endTime,
            isActive: true,
            entryFee: _entryFee,
            rewardPool: _rewardPool
        });

        emit LeagueCreated(leagueCount, _name, msg.sender);
    }

    // Join an existing league
    function joinLeague(uint256 _leagueId) public payable {
        League storage league = leagues[_leagueId];
        require(league.isActive, "League is not active");
        require(msg.value == league.entryFee, "Incorrect entry fee");

        leagueParticipants[_leagueId].push(
            Participant({user: msg.sender, points: 0})
        );

        emit JoinedLeague(_leagueId, msg.sender);
    }

    // Update points for a participant in a league
    function updatePoints(
        uint256 _leagueId,
        address _participant,
        uint256 _points
    ) public {
        League storage league = leagues[_leagueId];
        require(league.creator == msg.sender, "Only creator can update points");

        Participant[] storage participants = leagueParticipants[_leagueId];
        for (uint256 i = 0; i < participants.length; i++) {
            if (participants[i].user == _participant) {
                participants[i].points += _points;
                emit PointsUpdated(_leagueId, _participant, participants[i].points);
                return;
            }
        }
    }

    // Distribute rewards based on rankings
    function distributeRewards(uint256 _leagueId) public {
        League storage league = leagues[_leagueId];
        require(league.isActive, "League is already closed");
        require(block.timestamp > league.endTime, "League is still ongoing");

        Participant[] storage participants = leagueParticipants[_leagueId];
        // Sort participants by points
        quicksort(participants, 0, int256(participants.length - 1));

        // Distribute rewards
        uint256 totalReward = league.rewardPool;
        uint256 rewardPerParticipant = totalReward / participants.length;

        for (uint256 i = 0; i < participants.length; i++) {
            lingToken.transferFrom(treasury, participants[i].user, rewardPerParticipant);
        }

        league.isActive = false;
        emit RewardsDistributed(_leagueId);
    }

    // Quicksort and partition functions
    function quicksort(Participant[] storage arr, int256 left, int256 right) internal {
        if (left < right) {
            int256 pivotIndex = partition(arr, left, right);
            quicksort(arr, left, pivotIndex - 1);
            quicksort(arr, pivotIndex + 1, right);
        }
    }

    function partition(Participant[] storage arr, int256 left, int256 right) internal returns (int256) {
        uint256 pivot = arr[uint256(right)].points;
        int256 i = left - 1;

        for (int256 j = left; j < right; j++) {
            if (arr[uint256(j)].points > pivot) {
                i++;
                (arr[uint256(i)], arr[uint256(j)]) = (arr[uint256(j)], arr[uint256(i)]);
            }
        }

        i++;
        (arr[uint256(i)], arr[uint256(right)]) = (arr[uint256(right)], arr[uint256(i)]);
        return i;
    }
}