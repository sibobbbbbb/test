import CommunityEvent from '../models/CommunityEvent.js';
import { successResponse, errorResponse } from '../utils/response.js';
// import User from '../models/user.js';

export const getCommunityEvents = async (req, res) => {
  try {
    let events;
    // Jika user memiliki role "Buddy", hanya tampilkan event dengan accessLevel "Member and Buddy"
    if (req.user && req.user.access === 'Buddy') {
      events = await CommunityEvent.find({ accessLevel: 'Member and Buddy' });
    } else {
      events = await CommunityEvent.find();
    }
    return successResponse(res, 200, 'Success', events);
  } catch (error) {
    return errorResponse(res, 500, error.message);
  }
};

export const getCommunityEvent = async (req, res) => {
  try {
    const event = await CommunityEvent.findById(req.params.id)
      .populate('createdBy', 'name email')
      .populate('attendees.userId', 'name email');
    if (!event) {
      return res.status(404).json({ message: "Community Event not found" });
    }
    return successResponse(res, event);
  } catch (error) {
    return errorResponse(res, error.message);
  }
};

export const createCommunityEvent = async (req, res) => {
  try {
    const { title, description, date, location, accessLevel, capacity } = req.body;
    const newEvent = new CommunityEvent({
      title,
      description,
      date,
      location,
      accessLevel,
      capacity,
      createdBy: req.user._id
    });
    const savedEvent = await newEvent.save();
    
    // Opsional: Notifikasi via email atau push kepada user jika diperlukan
    // sendEmail({ to: ..., subject: ..., text: ... });
    // sendNotification({ userId: ..., message: ... });
    
    return res.status(201).json({ success: true, data: savedEvent });
  } catch (error) {
    return errorResponse(res, error.message);
  }
};

export const updateCommunityEvent = async (req, res) => {
  try {
    const updatedEvent = await CommunityEvent.findByIdAndUpdate(
      req.params.id,
      req.body,
      { new: true, runValidators: true }
    );
    if (!updatedEvent) {
      return res.status(404).json({ message: "Community Event not found" });
    }
    return successResponse(res, updatedEvent);
  } catch (error) {
    return errorResponse(res, error.message);
  }
};

export const deleteCommunityEvent = async (req, res) => {
  try {
    const deletedEvent = await CommunityEvent.findByIdAndDelete(req.params.id);
    if (!deletedEvent) {
      return res.status(404).json({ message: "Community Event not found" });
    }
    return successResponse(res, deletedEvent);
  } catch (error) {
    return errorResponse(res, error.message);
  }
};

export const rsvpCommunityEvent = async (req, res) => {
  try {
    const event = await CommunityEvent.findById(req.params.id);
    if (!event) {
      return res.status(404).json({ message: "Community Event not found" });
    }
    // Cek apakah user sudah RSVP
    const alreadyRSVP = event.attendees.some(att => att.userId.toString() === req.user._id.toString());
    if (alreadyRSVP) {
      return res.status(400).json({ message: "You have already RSVPed for this event" });
    }
    // Tambahkan RSVP user
    event.attendees.push({ userId: req.user._id, rsvp: true, attended: false });
    await event.save();
    return successResponse(res, event);
  } catch (error) {
    return errorResponse(res, error.message);
  }
};

export const cancelRsvpCommunityEvent = async (req, res) => {
  try {
    const event = await CommunityEvent.findById(req.params.id);
    if (!event) {
      return res.status(404).json({ message: "Community Event not found" });
    }
    // Hapus RSVP user
    event.attendees = event.attendees.filter(att => att.userId.toString() !== req.user._id.toString());
    await event.save();
    return successResponse(res, event);
  } catch (error) {
    return errorResponse(res, error.message);
  }
};

export const markAttendance = async (req, res) => {
  try {
    const { id, userId } = req.params;
    const event = await CommunityEvent.findById(id);
    if (!event) {
      return res.status(404).json({ message: "Community Event not found" });
    }
    // Cari attendee berdasarkan userId
    const attendee = event.attendees.find(att => att.userId.toString() === userId);
    if (!attendee) {
      return res.status(404).json({ message: "Attendee not found" });
    }
    attendee.attended = true;
    await event.save();
    return successResponse(res, event);
  } catch (error) {
    return errorResponse(res, error.message);
  }
};