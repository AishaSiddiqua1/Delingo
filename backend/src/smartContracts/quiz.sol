// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract QuizCompetition {
    struct Quiz {
        address creator;               // The address of the quiz creator
        string name;                   // The name of the quiz
        uint256 rewardPool;            // Total reward pool in Ling tokens for the quiz
        uint256 questionCount;         // Number of questions in the quiz
        uint256 startTime;             // Start time of the quiz
        uint256 endTime;               // End time of the quiz
        bool isActive;                 // Whether the quiz is still active
    }

    struct Participant {
        address user;                 // The user participating in the quiz
        uint256 score;                 // Score of the participant in the quiz
    }

    mapping(uint256 => Quiz) public quizzes;        // Mapping quizId to Quiz struct
    mapping(uint256 => Participant[]) public quizParticipants;  // Mapping quizId to participants array

    uint256 public quizCount = 0;                    // Counter for quizzes
    IERC20 public lingToken;                         // Ling token reference
    address public treasury;                         // Treasury address for transferring Ling tokens

    event QuizCreated(uint256 quizId, string name, address indexed creator);
    event QuizSolved(uint256 quizId, address indexed participant, uint256 score);
    event RewardsDistributed(uint256 quizId);

    constructor(address _lingToken, address _treasury) {
        lingToken = IERC20(_lingToken);
        treasury = _treasury;
    }

    // Modifier to ensure only the creator can perform certain actions
    modifier onlyCreator(uint256 _quizId) {
        require(quizzes[_quizId].creator == msg.sender, "Only the creator can perform this action");
        _;
    }

    // Create a new quiz
    function createQuiz(
        string memory _name,
        uint256 _rewardPool,
        uint256 _questionCount,
        uint256 _startTime,
        uint256 _endTime
    ) public {
        require(_startTime < _endTime, "Start time must be before end time");
        require(lingToken.allowance(treasury, address(this)) >= _rewardPool, "Insufficient reward pool allowance");

        quizCount++;
        quizzes[quizCount] = Quiz({
            creator: msg.sender,
            name: _name,
            rewardPool: _rewardPool,
            questionCount: _questionCount,
            startTime: _startTime,
            endTime: _endTime,
            isActive: true
        });

        emit QuizCreated(quizCount, _name, msg.sender);
    }

    // Solve a quiz (participants can answer questions)
    function solveQuiz(uint256 _quizId, uint256 _score) public {
        Quiz storage quiz = quizzes[_quizId];
        require(quiz.isActive, "Quiz is not active");
        require(block.timestamp >= quiz.startTime && block.timestamp <= quiz.endTime, "Quiz is not within the active period");

        // Add participant's score
        quizParticipants[_quizId].push(Participant({
            user: msg.sender,
            score: _score
        }));

        emit QuizSolved(_quizId, msg.sender, _score);
    }

    // Distribute rewards based on participant scores (highest score gets the most reward)
    function distributeRewards(uint256 _quizId) public onlyCreator(_quizId) {
        Quiz storage quiz = quizzes[_quizId];
        require(quiz.isActive, "Quiz is already closed");
        require(block.timestamp > quiz.endTime, "Quiz is still ongoing");

        Participant[] storage participants = quizParticipants[_quizId];
        uint256 totalReward = quiz.rewardPool;
        uint256 totalScore = 0;

        // Calculate total score
        for (uint256 i = 0; i < participants.length; i++) {
            totalScore += participants[i].score;
        }

        // If total score is zero, don't distribute rewards
        if (totalScore == 0) {
            return;
        }

        // Distribute rewards based on score percentage
        for (uint256 i = 0; i < participants.length; i++) {
            uint256 participantReward = (participants[i].score * totalReward) / totalScore;
            lingToken.transferFrom(treasury, participants[i].user, participantReward);
        }

        quiz.isActive = false;
        emit RewardsDistributed(_quizId);
    }

    // View function to get the details of a quiz
    function getQuizDetails(uint256 _quizId) public view returns (Quiz memory) {
        return quizzes[_quizId];
    }

    // View function to get the list of participants in a quiz
    function getQuizParticipants(uint256 _quizId) public view returns (Participant[] memory) {
        return quizParticipants[_quizId];
    }
}
