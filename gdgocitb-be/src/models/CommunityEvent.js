import mongoose from 'mongoose';

const CommunityEventSchema = new mongoose.Schema({
  title: {
    type: String,
    required: [true, 'Event title is required'],
    trim: true
  },
  subtitle: {
    type: String,
    trim: true
  },
  category: {
    type: String,
    default: 'Google Developer Groups',
    enum: ['Google Developer Groups']
  },
  description: {
    type: String,
    trim: true
  },
  location: {
    type: String,
    trim: true
  },
  date: {
    type: Date,
    required: [true, 'Event date is required']
  },
  time: {
    start: {
      type: String,
      required: [true, 'Event start time is required']
    },
    end: {
      type: String,
      required: [true, 'Event end time is required']
    }
  },
  accessLevel: {
    type: String,
    enum: ['Member', 'Member and Buddy'],
    default: 'Member'
  },
  capacity: {
    type: Number
  },
  attendees: [{
    userId: {
      type: mongoose.Schema.Types.ObjectId,
      ref: 'User'
    },
    rsvp: {
      type: Boolean,
      default: false
    },
    attended: {
      type: Boolean,
      default: false
    }
  }],
  createdAt: {
    type: Date,
    default: Date.now
  },
  updatedAt: {
    type: Date,
    default: Date.now
  }
}, {
  timestamps: true
});

const CommunityEvent = mongoose.model('CommunityEvent', CommunityEventSchema);

export default CommunityEvent;